package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"

	"github.com/gin-gonic/gin"
)

const maxStudentEvidenceSize = 5 << 20

var studentApplicationDir = "uploads/student-applications"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StudentResolveRequest struct {
	Approved bool `json:"approved" binding:"required"`
}

type AuthController struct {
	repo *repositories.UserRepository
}

func NewAuthController(repo *repositories.UserRepository) *AuthController {
	return &AuthController{repo: repo}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Username = strings.TrimSpace(request.Username)
	if len(request.Username) < 3 || len(request.Username) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be 3-32 characters"})
		return
	}
	if len(request.Password) < 4 || len(request.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be 4-72 characters"})
		return
	}

	user, err := ctrl.repo.CreateUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registration successful", "user": user})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, found := ctrl.repo.GetUserByUsername(request.Username)
	if !found || user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}
	token := hex.EncodeToString(tokenBytes)
	ctrl.repo.SaveSession(token, user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token, "user": user})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	token, _ := c.Get("authToken")
	ctrl.repo.DeleteSession(token.(string))
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (ctrl *AuthController) Me(c *gin.Context) {
	user, _ := CurrentUser(c)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *AuthController) SubmitStudentApplication(c *gin.Context) {
	user, _ := CurrentUser(c)
	if user.IsAdmin == models.AdminTrue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin cannot apply for student status"})
		return
	}
	if user.IsStudent == models.StudentTrue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is already a verified student"})
		return
	}
	if user.IsStudent == models.StudentPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student application is already pending"})
		return
	}
	file, err := c.FormFile("evidence")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "evidence file is required"})
		return
	}
	if file.Size <= 0 || file.Size > maxStudentEvidenceSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "evidence must be at most 5 MB"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read evidence"})
		return
	}
	defer opened.Close()

	data, err := io.ReadAll(io.LimitReader(opened, maxStudentEvidenceSize+1))
	if err != nil || len(data) > maxStudentEvidenceSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read evidence"})
		return
	}
	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "evidence must be JPG, PNG, or PDF"})
		return
	}

	if err := os.MkdirAll(studentApplicationDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store evidence"})
		return
	}
	extension := map[string]string{
		"image/jpeg":      ".jpg",
		"image/png":       ".png",
		"application/pdf": ".pdf",
	}[contentType]
	evidencePath := filepath.Join(studentApplicationDir, strconv.Itoa(user.ID)+extension)
	if err := os.WriteFile(evidencePath, data, 0600); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store evidence"})
		return
	}

	updatedUser, err := ctrl.repo.SubmitStudentApplication(user.ID, filepath.Base(file.Filename), contentType, evidencePath)
	if err != nil {
		_ = os.Remove(evidencePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "student application submitted", "data": updatedUser})
}

func (ctrl *AuthController) ListStudentApplications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": ctrl.repo.GetPendingStudentApplications()})
}

func (ctrl *AuthController) StudentEvidence(c *gin.Context) {
	userID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	application, found := ctrl.repo.GetStudentApplication(userID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "student application not found"})
		return
	}

	disposition := mime.FormatMediaType("inline", map[string]string{"filename": application.Filename})
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", application.ContentType)
	if strings.Contains(application.EvidencePath, "..") {
    	c.JSON(http.StatusForbidden, gin.H{"error": "invalid evidence path"})
    	return
	}
	c.File(application.EvidencePath)
}

func (ctrl *AuthController) ResolveStudentApplication(c *gin.Context) {
	var approved StudentResolveRequest

	if err := c.ShouldBindJSON(&approved); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approval status is required"})
		return
	}

	userID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	application, found := ctrl.repo.GetStudentApplication(userID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "student application not found"})
		return
	}
	user, err := ctrl.repo.ResolveStudentApplication(userID, approved.Approved)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	_ = os.Remove(application.EvidencePath)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *AuthController) RequireAuth(c *gin.Context) {
	user, token, ok := ctrl.authenticate(c)
	if !ok {
		return
	}
	c.Set("currentUser", user)
	c.Set("authToken", token)
	c.Next()
}

func (ctrl *AuthController) RequireAdmin(c *gin.Context) {
	user, token, ok := ctrl.authenticate(c)
	if !ok {
		return
	}
	if user.IsAdmin != models.AdminTrue {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	c.Set("currentUser", user)
	c.Set("authToken", token)
	c.Next()
}

func (ctrl *AuthController) RequireNonAdmin(c *gin.Context) {
	user, _ := CurrentUser(c)
	if user.IsAdmin == models.AdminTrue {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin users cannot make bookings or payments"})
		return
	}
	c.Next()
}

func (ctrl *AuthController) authenticate(c *gin.Context) (models.User, string, bool) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return models.User{}, "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	user, ok := ctrl.repo.GetUserByToken(token)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return models.User{}, "", false
	}
	return user, token, true
}

func CurrentUser(c *gin.Context) (models.User, bool) {
	value, ok := c.Get("currentUser")
	if !ok {
		return models.User{}, false
	}
	user, ok := value.(models.User)
	return user, ok
}

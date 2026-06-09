package controllers

import (
	"net/http"

	"gin-M-TIX/config"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   "demo-token",
		"user": gin.H{
			"id":       1,
			"username": request.Username,
		},
	})
}

func Logout(db *config.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		db.Reset()
		c.JSON(http.StatusOK, gin.H{"message": "logout successful, database reset"})
	}
}

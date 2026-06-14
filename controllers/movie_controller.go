package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"

	"github.com/gin-gonic/gin"
)

const maxPosterSize = 5 << 20

type MovieController struct {
	repo *repositories.MovieRepository
}

func NewMovieController(repo *repositories.MovieRepository) *MovieController {
	return &MovieController{repo: repo}
}

func (ctrl *MovieController) GetMovies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": ctrl.repo.GetAll()})
}

func (ctrl *MovieController) CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posterURL, err := savePoster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie.PosterURL = posterURL
	c.JSON(http.StatusCreated, gin.H{"data": ctrl.repo.Create(movie)})
}

func (ctrl *MovieController) UpdateMovie(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var movie models.Movie
	if err := c.ShouldBind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	current, found := ctrl.repo.GetByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	movie.PosterURL = current.PosterURL
	if _, err := c.FormFile("poster"); err == nil {
		posterURL, err := savePoster(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		movie.PosterURL = posterURL
	}
	updatedMovie, err := ctrl.repo.Update(id, movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if movie.PosterURL != current.PosterURL {
		removeUploadedPoster(current.PosterURL)
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedMovie})
}

func (ctrl *MovieController) DeleteMovie(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	movie, found := ctrl.repo.GetByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	if err := ctrl.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	removeUploadedPoster(movie.PosterURL)

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}

func savePoster(c *gin.Context) (string, error) {
	file, err := c.FormFile("poster")
	if err != nil {
		return "", fmt.Errorf("poster file is required")
	}
	if file.Size <= 0 || file.Size > maxPosterSize {
		return "", fmt.Errorf("poster must be at most 5 MB")
	}

	opened, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to read poster")
	}
	defer opened.Close()

	data, err := io.ReadAll(io.LimitReader(opened, maxPosterSize+1))
	if err != nil || len(data) > maxPosterSize {
		return "", fmt.Errorf("failed to read poster")
	}
	contentType := http.DetectContentType(data)
	extensions := map[string]string{"image/jpeg": ".jpg", "image/png": ".png"}
	extension, ok := extensions[contentType]
	if !ok {
		return "", fmt.Errorf("poster must be JPG or PNG")
	}

	if err := os.MkdirAll("public/poster", 0755); err != nil {
		return "", fmt.Errorf("failed to store poster")
	}
	target, err := os.CreateTemp("public/poster", "movie-*"+extension)
	if err != nil {
		return "", fmt.Errorf("failed to store poster")
	}
	path := target.Name()
	if _, err := target.Write(data); err != nil {
		target.Close()
		os.Remove(path)
		return "", fmt.Errorf("failed to store poster")
	}
	if err := target.Close(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to store poster")
	}
	basePath := filepath.Base(path)
	if strings.Contains(basePath, "..") {
    	os.Remove(path)
    	return "", fmt.Errorf("invalid poster path")
	}
	return "/ui/poster/" + basePath, nil
}

func removeUploadedPoster(url string) {
	name := filepath.Base(url)
    if !strings.HasPrefix(name, "movie-") {
        return
    }
    target := filepath.Join("public/poster", name)
    if rel, err := filepath.Rel("public/poster", target); err != nil || strings.Contains(rel, "..") {
        return
    }
    _ = os.Remove(target)
}

func parseIDParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

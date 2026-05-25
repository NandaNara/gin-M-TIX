package controllers

import (
	"net/http"
	"strconv"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"

	"github.com/gin-gonic/gin"
)

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
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ctrl.repo.Create(movie)})
}

func (ctrl *MovieController) UpdateMovie(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMovie, err := ctrl.repo.Update(id, movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedMovie})
}

func (ctrl *MovieController) DeleteMovie(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := ctrl.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}

func parseIDParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

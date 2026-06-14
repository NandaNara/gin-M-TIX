package controllers

import (
	"net/http"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"

	"github.com/gin-gonic/gin"
)

type StudioController struct {
	repo *repositories.StudioRepository
}

func NewStudioController(repo *repositories.StudioRepository) *StudioController {
	return &StudioController{repo: repo}
}

func (ctrl *StudioController) GetStudios(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": ctrl.repo.GetStudios()})
}

func (ctrl *StudioController) CreateStudio(c *gin.Context) {
	var studio models.Studio
	if err := c.ShouldBindJSON(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdStudio, err := ctrl.repo.CreateStudio(studio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": createdStudio})
}

func (ctrl *StudioController) UpdateStudio(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var studio models.Studio
	if err := c.ShouldBindJSON(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := ctrl.repo.UpdateStudio(id, studio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (ctrl *StudioController) DeleteStudio(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := ctrl.repo.DeleteStudio(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "studio deleted"})
}

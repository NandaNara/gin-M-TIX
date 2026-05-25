package controllers

import (
	"net/http"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"
	"gin-M-TIX/services"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	repo    *repositories.ScheduleRepository
	pricing *services.PricingService
}

type ScheduleResponse struct {
	models.Schedule
	SeatPrice       float64 `json:"seat_price"`
	PricingStrategy string  `json:"pricing_strategy"`
}

func NewScheduleController(repo *repositories.ScheduleRepository, pricing *services.PricingService) *ScheduleController {
	return &ScheduleController{repo: repo, pricing: pricing}
}

func (ctrl *ScheduleController) GetSchedules(c *gin.Context) {
	schedules := ctrl.repo.GetAll()
	responses := make([]ScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		responses = append(responses, ctrl.toResponse(schedule))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (ctrl *ScheduleController) CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSchedule, err := ctrl.repo.Create(schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ctrl.toResponse(createdSchedule)})
}

func (ctrl *ScheduleController) GetScheduleSeats(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	seats, err := ctrl.repo.GetSeatsByScheduleID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seats})
}

func (ctrl *ScheduleController) toResponse(schedule models.Schedule) ScheduleResponse {
	seatPrice, strategyName := ctrl.pricing.CalculateSeatPrice(schedule)
	return ScheduleResponse{
		Schedule:        schedule,
		SeatPrice:       seatPrice,
		PricingStrategy: strategyName,
	}
}

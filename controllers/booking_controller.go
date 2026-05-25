package controllers

import (
	"net/http"

	bookingfacade "gin-M-TIX/patterns/facade"
	"gin-M-TIX/services"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	facade *bookingfacade.BookingFacade
}

func NewBookingController(facade *bookingfacade.BookingFacade) *BookingController {
	return &BookingController{facade: facade}
}

func (ctrl *BookingController) CreateBooking(c *gin.Context) {
	var request services.CreateBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := ctrl.facade.CreateBooking(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

func (ctrl *BookingController) GetBooking(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	booking, found := ctrl.facade.GetBooking(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (ctrl *BookingController) GetUserBookings(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ctrl.facade.GetUserBookings(id)})
}

func (ctrl *BookingController) Pay(c *gin.Context) {
	var request services.PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, booking, err := ctrl.facade.Pay(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"payment": payment,
			"booking": booking,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"payment": payment,
			"booking": booking,
		},
	})
}

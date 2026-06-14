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
	user, _ := CurrentUser(c)
	request.UserID = user.ID

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

	user, _ := CurrentUser(c)
	booking, found := ctrl.facade.GetBookingForUser(id, user.ID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (ctrl *BookingController) GetUserBookings(c *gin.Context) {
	user, _ := CurrentUser(c)
	c.JSON(http.StatusOK, gin.H{"data": ctrl.facade.GetUserBookings(user.ID)})
}

func (ctrl *BookingController) CancelBooking(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	user, _ := CurrentUser(c)
	booking, err := ctrl.facade.CancelBookingForUser(id, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking, "message": "Booking canceled successfully"})
}

func (ctrl *BookingController) Pay(c *gin.Context) {
	var request services.PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := CurrentUser(c)
	payment, booking, err := ctrl.facade.Pay(user.ID, request)
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

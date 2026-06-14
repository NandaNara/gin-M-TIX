package facade

import (
	"gin-M-TIX/models"
	"gin-M-TIX/services"
)

type BookingFacade struct {
	bookingService *services.BookingService
	paymentService *services.PaymentService
}

func NewBookingFacade(
	bookingService *services.BookingService,
	paymentService *services.PaymentService,
) *BookingFacade {
	return &BookingFacade{
		bookingService: bookingService,
		paymentService: paymentService,
	}
}

func (f *BookingFacade) CreateBooking(request services.CreateBookingRequest) (models.Booking, error) {
	return f.bookingService.CreateBooking(request)
}

func (f *BookingFacade) GetBookingForUser(id, userID int) (models.Booking, bool) {
	return f.bookingService.GetBookingForUser(id, userID)
}

func (f *BookingFacade) GetUserBookings(userID int) []models.Booking {
	return f.bookingService.GetUserBookings(userID)
}

func (f *BookingFacade) CancelBookingForUser(id, userID int) (models.Booking, error) {
	return f.bookingService.CancelBookingForUser(id, userID)
}

func (f *BookingFacade) Pay(userID int, request services.PaymentRequest) (models.Payment, models.Booking, error) {
	return f.paymentService.Pay(userID, request)
}

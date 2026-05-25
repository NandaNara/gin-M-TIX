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

func (f *BookingFacade) GetBooking(id int) (models.Booking, bool) {
	return f.bookingService.GetBooking(id)
}

func (f *BookingFacade) GetUserBookings(userID int) []models.Booking {
	return f.bookingService.GetUserBookings(userID)
}

func (f *BookingFacade) Pay(request services.PaymentRequest) (models.Payment, models.Booking, error) {
	return f.paymentService.Pay(request)
}

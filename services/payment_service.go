package services

import (
	"errors"

	"gin-M-TIX/models"
	"gin-M-TIX/repositories"
)

type PaymentRequest struct {
	BookingID int     `json:"booking_id" binding:"required"`
	Method    string  `json:"method" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

type PaymentService struct {
	bookingRepo *repositories.BookingRepository
}

func NewPaymentService(bookingRepo *repositories.BookingRepository) *PaymentService {
	return &PaymentService{bookingRepo: bookingRepo}
}

func (s *PaymentService) Pay(request PaymentRequest) (models.Payment, models.Booking, error) {
	booking, ok := s.bookingRepo.GetByID(request.BookingID)
	if !ok {
		return models.Payment{}, models.Booking{}, errors.New("booking not found")
	}
	if booking.Status == models.BookingPaid {
		return models.Payment{}, booking, errors.New("booking already paid")
	}

	status := models.PaymentSuccess
	if request.Amount < booking.TotalPrice {
		status = models.PaymentFailed
		payment := s.bookingRepo.CreatePayment(models.Payment{
			BookingID: request.BookingID,
			Method:    request.Method,
			Amount:    request.Amount,
			Status:    status,
		})
		return payment, booking, errors.New("payment amount is less than total price")
	}

	payment := s.bookingRepo.CreatePayment(models.Payment{
		BookingID: request.BookingID,
		Method:    request.Method,
		Amount:    request.Amount,
		Status:    status,
	})
	paidBooking, err := s.bookingRepo.MarkPaid(request.BookingID)
	if err != nil {
		return payment, booking, err
	}
	return payment, paidBooking, nil
}

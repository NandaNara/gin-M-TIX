package services

import (
	"errors"
	"fmt"

	"gin-M-TIX/models"
	ticketfactory "gin-M-TIX/patterns/factory"
	"gin-M-TIX/repositories"
)

type CreateBookingRequest struct {
	UserID     int   `json:"-"`
	ScheduleID int   `json:"schedule_id" binding:"required"`
	SeatIDs    []int `json:"seat_ids" binding:"required"`
}

type BookingService struct {
	bookingRepo  *repositories.BookingRepository
	scheduleRepo *repositories.ScheduleRepository
	pricing      *PricingService
}

func NewBookingService(
	bookingRepo *repositories.BookingRepository,
	scheduleRepo *repositories.ScheduleRepository,
	pricing *PricingService,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
		pricing:      pricing,
	}
}

func (s *BookingService) CreateBooking(request CreateBookingRequest) (models.Booking, error) {
	if request.UserID <= 0 {
		return models.Booking{}, errors.New("user_id must be greater than zero")
	}
	if s.bookingRepo.IsUserAdmin(request.UserID) {
		return models.Booking{}, errors.New("admin users cannot create bookings")
	}
	if len(request.SeatIDs) == 0 {
		return models.Booking{}, errors.New("seat_ids cannot be empty")
	}

	schedule, ok := s.scheduleRepo.GetByID(request.ScheduleID)
	if !ok {
		return models.Booking{}, errors.New("schedule not found")
	}

	if !s.bookingRepo.IsBookedDayValid(schedule) {
		return models.Booking{}, errors.New("booking must be made at least 3 days in advance and before the schedule start time")
	}

	availableSeats, err := s.scheduleRepo.GetSeatsByScheduleID(request.ScheduleID)
	if err != nil {
		return models.Booking{}, err
	}

	seatMap := make(map[int]models.Seat, len(availableSeats))
	for _, seat := range availableSeats {
		seatMap[seat.ID] = seat
	}

	uniqueSeatIDs := make(map[int]bool, len(request.SeatIDs))
	selectedSeats := make([]models.Seat, 0, len(request.SeatIDs))
	for _, seatID := range request.SeatIDs {
		if uniqueSeatIDs[seatID] {
			return models.Booking{}, fmt.Errorf("duplicate seat id: %d", seatID)
		}
		uniqueSeatIDs[seatID] = true

		seat, ok := seatMap[seatID]
		if !ok {
			return models.Booking{}, fmt.Errorf("seat id %d does not belong to schedule studio", seatID)
		}
		if seat.Status == models.SeatBooked {
			return models.Booking{}, fmt.Errorf("seat %s is already booked", seat.Code)
		}
		selectedSeats = append(selectedSeats, seat)
	}

	baseSeatPrice, _ := s.pricing.CalculateSeatPrice(schedule)
	isStudent := s.bookingRepo.IsUserStudent(request.UserID)

	tickets := make([]models.Ticket, 0, len(selectedSeats))
	totalPrice := 0.0
	for _, seat := range selectedSeats {
		factory := ticketfactory.NewTicketFactory(bool(seat.IsVIP), isStudent)
		ticket := factory.CreateTicket(schedule.ID, seat, baseSeatPrice)
		tickets = append(tickets, ticket)
		totalPrice += ticket.Price
	}

	booking := models.Booking{
		UserID:     request.UserID,
		ScheduleID: request.ScheduleID,
		SeatIDs:    request.SeatIDs,
		Tickets:    tickets,
		TotalPrice: totalPrice,
	}

	return s.bookingRepo.Create(booking)
}

func (s *BookingService) GetBookingForUser(id, userID int) (models.Booking, bool) {
	booking, ok := s.bookingRepo.GetByID(id)
	return booking, ok && booking.UserID == userID
}

func (s *BookingService) GetUserBookings(userID int) []models.Booking {
	return s.bookingRepo.GetByUserID(userID)
}

func (s *BookingService) CancelBookingForUser(id, userID int) (models.Booking, error) {
	booking, ok := s.bookingRepo.GetByID(id)
	if !ok || booking.UserID != userID {
		return models.Booking{}, errors.New("booking not found")
	}
	return s.bookingRepo.Cancel(id)
}

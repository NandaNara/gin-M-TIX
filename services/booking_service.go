package services

import (
	"errors"
	"fmt"

	"gin-M-TIX/models"
	ticketfactory "gin-M-TIX/patterns/factory"
	"gin-M-TIX/repositories"
)

type CreateBookingRequest struct {
	UserID     int               `json:"user_id" binding:"required"`
	ScheduleID int               `json:"schedule_id" binding:"required"`
	SeatIDs    []int             `json:"seat_ids" binding:"required"`
	TicketType models.TicketType `json:"ticket_type"`
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
	if len(request.SeatIDs) == 0 {
		return models.Booking{}, errors.New("seat_ids cannot be empty")
	}
	if request.TicketType == "" {
		request.TicketType = models.TicketRegular
	}

	schedule, ok := s.scheduleRepo.GetByID(request.ScheduleID)
	if !ok {
		return models.Booking{}, errors.New("schedule not found")
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
	factory, err := ticketfactory.NewTicketFactory(request.TicketType)
	if err != nil {
		return models.Booking{}, err
	}

	tickets := make([]models.Ticket, 0, len(selectedSeats))
	totalPrice := 0.0
	for _, seat := range selectedSeats {
		ticket := factory.CreateTicket(schedule.ID, seat, baseSeatPrice)
		tickets = append(tickets, ticket)
		totalPrice += ticket.Price
	}

	booking := models.Booking{
		UserID:     request.UserID,
		ScheduleID: request.ScheduleID,
		SeatIDs:    request.SeatIDs,
		TicketType: request.TicketType,
		Tickets:    tickets,
		TotalPrice: totalPrice,
	}

	return s.bookingRepo.Create(booking)
}

func (s *BookingService) GetBooking(id int) (models.Booking, bool) {
	return s.bookingRepo.GetByID(id)
}

func (s *BookingService) GetUserBookings(userID int) []models.Booking {
	return s.bookingRepo.GetByUserID(userID)
}

func (s *BookingService) CancelBooking(id int) (models.Booking, error) {
	return s.bookingRepo.Cancel(id)
}

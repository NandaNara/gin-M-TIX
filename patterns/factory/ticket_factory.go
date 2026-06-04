package factory

import (
	"errors"

	"gin-M-TIX/models"
)

type TicketFactory interface {
	CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket
}

type RegularTicketFactory struct{}

func (RegularTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{
		ScheduleID: scheduleID,
		SeatID:     seat.ID,
		SeatCode:   seat.Code,
		Type:       models.TicketRegular,
		Price:      basePrice,
	}
}

type VIPTicketFactory struct{}

func (VIPTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{
		ScheduleID: scheduleID,
		SeatID:     seat.ID,
		SeatCode:   seat.Code,
		Type:       models.TicketVIP,
		Price:      basePrice * 1.5,
	}
}

type StudentTicketFactory struct{}

func (StudentTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{
		ScheduleID: scheduleID,
		SeatID:     seat.ID,
		SeatCode:   seat.Code,
		Type:       models.TicketStudent,
		Price:      basePrice * 0.8, // 20% discount for students
	}
}

func NewTicketFactory(ticketType models.TicketType) (TicketFactory, error) {
	switch ticketType {
	case "", models.TicketRegular:
		return RegularTicketFactory{}, nil
	case models.TicketStudent:
		return StudentTicketFactory{}, nil
	case models.TicketVIP:
		return VIPTicketFactory{}, nil
	default:
		return nil, errors.New("unsupported ticket type")
	}
}

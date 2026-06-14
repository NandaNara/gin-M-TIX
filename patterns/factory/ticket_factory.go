package factory

import "gin-M-TIX/models"

type TicketFactory interface {
	CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket
}

type RegularTicketFactory struct{}

func (RegularTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{ScheduleID: scheduleID, SeatID: seat.ID, SeatCode: seat.Code, Type: models.TicketRegular, Price: basePrice}
}

type VIPTicketFactory struct{}

func (VIPTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{ScheduleID: scheduleID, SeatID: seat.ID, SeatCode: seat.Code, Type: models.TicketVIP, Price: basePrice * 1.5}
}

type RegularStudentTicketFactory struct{}

func (RegularStudentTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{ScheduleID: scheduleID, SeatID: seat.ID, SeatCode: seat.Code, Type: models.TicketRegularStudent, Price: basePrice * 0.8}
}

type VIPStudentTicketFactory struct{}

func (VIPStudentTicketFactory) CreateTicket(scheduleID int, seat models.Seat, basePrice float64) models.Ticket {
	return models.Ticket{ScheduleID: scheduleID, SeatID: seat.ID, SeatCode: seat.Code, Type: models.TicketVIPStudent, Price: basePrice * 1.5 * 0.8}
}

func NewTicketFactory(isVIP, isStudent bool) TicketFactory {
	if isVIP && isStudent {
		return VIPStudentTicketFactory{}
	}
	if isVIP {
		return VIPTicketFactory{}
	}
	if isStudent {
		return RegularStudentTicketFactory{}
	}
	return RegularTicketFactory{}
}

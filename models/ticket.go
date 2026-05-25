package models

import "time"

type TicketType string

const (
	TicketRegular TicketType = "regular"
	TicketVIP     TicketType = "vip"
)

type Ticket struct {
	ID         int        `json:"id"`
	BookingID  int        `json:"booking_id"`
	ScheduleID int        `json:"schedule_id"`
	SeatID     int        `json:"seat_id"`
	SeatCode   string     `json:"seat_code"`
	Type       TicketType `json:"type"`
	Price      float64    `json:"price"`
	CreatedAt  time.Time  `json:"created_at"`
}

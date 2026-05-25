package models

import "time"

type BookingStatus string

const (
	BookingPending  BookingStatus = "pending_payment"
	BookingPaid     BookingStatus = "paid"
	BookingCanceled BookingStatus = "canceled"
)

type Booking struct {
	ID         int           `json:"id"`
	UserID     int           `json:"user_id" binding:"required"`
	ScheduleID int           `json:"schedule_id" binding:"required"`
	SeatIDs    []int         `json:"seat_ids" binding:"required"`
	TicketType TicketType    `json:"ticket_type"`
	Tickets    []Ticket      `json:"tickets"`
	TotalPrice float64       `json:"total_price"`
	Status     BookingStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	PaidAt     *time.Time    `json:"paid_at,omitempty"`
}

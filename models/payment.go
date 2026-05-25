package models

import "time"

type PaymentStatus string

const (
	PaymentSuccess PaymentStatus = "success"
	PaymentFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID        int           `json:"id"`
	BookingID int           `json:"booking_id" binding:"required"`
	Method    string        `json:"method" binding:"required"`
	Amount    float64       `json:"amount" binding:"required"`
	Status    PaymentStatus `json:"status"`
	PaidAt    time.Time     `json:"paid_at"`
}

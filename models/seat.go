package models

type SeatStatus string

const (
	SeatAvailable SeatStatus = "available"
	SeatBooked    SeatStatus = "booked"
)

type Seat struct {
	ID       int        `json:"id"`
	StudioID int        `json:"studio_id"`
	Row      string     `json:"row"`
	Number   int        `json:"number"`
	Code     string     `json:"code"`
	Status   SeatStatus `json:"status"`
}

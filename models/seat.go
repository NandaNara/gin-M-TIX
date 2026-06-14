package models

type SeatStatus string

const (
	SeatAvailable SeatStatus = "available"
	SeatBooked    SeatStatus = "booked"
)

type SeatVIPStatus bool

const (
	RegularSeat SeatVIPStatus = false
	VIPSeat     SeatVIPStatus = true
)

type Seat struct {
	ID       int           `json:"id"`
	StudioID int           `json:"studio_id"`
	Row      string        `json:"row"`
	Number   int           `json:"number"`
	Code     string        `json:"code"`
	Status   SeatStatus    `json:"status"`
	IsVIP    SeatVIPStatus `json:"is_vip"`
}

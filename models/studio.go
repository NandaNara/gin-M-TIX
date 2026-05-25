package models

type Studio struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SeatRows    int    `json:"seat_rows"`
	SeatColumns int    `json:"seat_columns"`
}

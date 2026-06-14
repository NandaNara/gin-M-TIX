package models

type Studio struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	SeatRows    int    `json:"seat_rows" binding:"required,min=1,max=26"`
	SeatColumns int    `json:"seat_columns" binding:"required,min=1,max=50"`
}

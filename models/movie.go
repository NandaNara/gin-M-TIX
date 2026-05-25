package models

type Movie struct {
	ID              int    `json:"id"`
	Title           string `json:"title" binding:"required"`
	Genre           string `json:"genre" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required"`
	Rating          string `json:"rating"`
}

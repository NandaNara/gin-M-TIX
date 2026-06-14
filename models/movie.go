package models

type Movie struct {
	ID              int    `json:"id"`
	Title           string `json:"title" form:"title" binding:"required"`
	Genre           string `json:"genre" form:"genre" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" form:"duration_minutes" binding:"required,min=1"`
	Rating          string `json:"rating" form:"rating"`
	PosterURL       string `json:"poster_url"`
}

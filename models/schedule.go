package models

import "time"

type Schedule struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id" binding:"required"`
	StudioID  int       `json:"studio_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	BasePrice float64   `json:"base_price" binding:"required"`
	Movie     *Movie    `json:"movie,omitempty"`
	Studio    *Studio   `json:"studio,omitempty"`
}

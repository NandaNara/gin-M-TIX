package strategy

import "gin-M-TIX/models"

type WeekendPricing struct{}

func (WeekendPricing) Calculate(schedule models.Schedule) float64 {
	return schedule.BasePrice * 1.25
}

func (WeekendPricing) Name() string {
	return "weekend"
}

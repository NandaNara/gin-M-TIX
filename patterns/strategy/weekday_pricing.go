package strategy

import "gin-M-TIX/models"

type WeekdayPricing struct{}

func (WeekdayPricing) Calculate(schedule models.Schedule) float64 {
	return schedule.BasePrice
}

func (WeekdayPricing) Name() string {
	return "weekday"
}

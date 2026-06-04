package strategy

import "gin-M-TIX/models"

type MidnightPricing struct{}

func (MidnightPricing) Calculate(schedule models.Schedule) float64 {
	return schedule.BasePrice * 1.2
}

func (MidnightPricing) Name() string {
	return "midnight"
}

package strategy

import "gin-M-TIX/models"

type HolidayPricing struct{}

func (HolidayPricing) Calculate(schedule models.Schedule) float64 {
	return schedule.BasePrice * 1.5
}

func (HolidayPricing) Name() string {
	return "holiday"
}

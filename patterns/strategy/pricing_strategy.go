package strategy

import (
	"time"

	"gin-M-TIX/models"
)

type PricingStrategy interface {
	Calculate(schedule models.Schedule) float64
	Name() string
}

func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

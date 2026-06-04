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

func IsMidnight(t time.Time) bool {
	return t.Hour() >= 22 || t.Hour() <= 2
}

func IsHoliday(t time.Time) bool {
	// Simple mock for holiday (e.g. New Year or Christmas)
	return (t.Month() == time.January && t.Day() == 1) || (t.Month() == time.December && t.Day() == 25)
}

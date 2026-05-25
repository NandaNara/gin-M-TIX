package services

import (
	"gin-M-TIX/models"
	"gin-M-TIX/patterns/strategy"
)

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (s *PricingService) GetStrategy(schedule models.Schedule) strategy.PricingStrategy {
	if strategy.IsWeekend(schedule.StartTime) {
		return strategy.WeekendPricing{}
	}
	return strategy.WeekdayPricing{}
}

func (s *PricingService) CalculateSeatPrice(schedule models.Schedule) (float64, string) {
	pricingStrategy := s.GetStrategy(schedule)
	return pricingStrategy.Calculate(schedule), pricingStrategy.Name()
}

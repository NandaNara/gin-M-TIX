package repositories

import (
	"errors"
	"sort"

	"gin-M-TIX/config"
	"gin-M-TIX/models"
)

type ScheduleRepository struct {
	db *config.Database
}

func NewScheduleRepository(db *config.Database) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetAll() []models.Schedule {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	schedules := make([]models.Schedule, 0, len(r.db.Schedules))
	for _, schedule := range r.db.Schedules {
		schedules = append(schedules, r.enrichScheduleLocked(schedule))
	}
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].StartTime.Before(schedules[j].StartTime)
	})
	return schedules
}

func (r *ScheduleRepository) GetByID(id int) (models.Schedule, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	schedule, ok := r.db.Schedules[id]
	if !ok {
		return models.Schedule{}, false
	}
	return r.enrichScheduleLocked(schedule), true
}

func (r *ScheduleRepository) Create(schedule models.Schedule) (models.Schedule, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, ok := r.db.Movies[schedule.MovieID]; !ok {
		return models.Schedule{}, errors.New("movie not found")
	}
	if _, ok := r.db.Studios[schedule.StudioID]; !ok {
		return models.Schedule{}, errors.New("studio not found")
	}

	schedule.ID = r.db.NextIDs["schedules"]
	r.db.NextIDs["schedules"]++
	r.db.Schedules[schedule.ID] = schedule
	return r.enrichScheduleLocked(schedule), nil
}

func (r *ScheduleRepository) GetSeatsByScheduleID(scheduleID int) ([]models.Seat, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	schedule, ok := r.db.Schedules[scheduleID]
	if !ok {
		return nil, errors.New("schedule not found")
	}

	bookedSeatIDs := make(map[int]bool)
	for _, booking := range r.db.Bookings {
		if booking.ScheduleID != scheduleID || booking.Status == models.BookingCanceled {
			continue
		}
		for _, seatID := range booking.SeatIDs {
			bookedSeatIDs[seatID] = true
		}
	}

	seats := make([]models.Seat, 0)
	for _, seat := range r.db.Seats {
		if seat.StudioID != schedule.StudioID {
			continue
		}
		if bookedSeatIDs[seat.ID] {
			seat.Status = models.SeatBooked
		} else {
			seat.Status = models.SeatAvailable
		}
		seats = append(seats, seat)
	}
	sort.Slice(seats, func(i, j int) bool {
		return seats[i].ID < seats[j].ID
	})
	return seats, nil
}

func (r *ScheduleRepository) SeatBelongsToSchedule(seatID int, scheduleID int) bool {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	schedule, scheduleOK := r.db.Schedules[scheduleID]
	seat, seatOK := r.db.Seats[seatID]
	return scheduleOK && seatOK && seat.StudioID == schedule.StudioID
}

func (r *ScheduleRepository) enrichScheduleLocked(schedule models.Schedule) models.Schedule {
	if movie, ok := r.db.Movies[schedule.MovieID]; ok {
		schedule.Movie = &movie
	}
	if studio, ok := r.db.Studios[schedule.StudioID]; ok {
		schedule.Studio = &studio
	}
	return schedule
}

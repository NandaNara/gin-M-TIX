package repositories

import (
	"errors"
	"sort"
	"time"

	"gin-M-TIX/config"
	"gin-M-TIX/models"
)

type BookingRepository struct {
	db *config.Database
}

func NewBookingRepository(db *config.Database) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking models.Booking) (models.Booking, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if r.hasConflictingSeatsLocked(booking.ScheduleID, booking.SeatIDs, 0) {
		return models.Booking{}, errors.New("one or more seats already booked")
	}

	booking.ID = r.db.NextIDs["bookings"]
	r.db.NextIDs["bookings"]++
	booking.CreatedAt = time.Now()
	booking.Status = models.BookingPending

	for index := range booking.Tickets {
		booking.Tickets[index].ID = r.db.NextIDs["tickets"]
		r.db.NextIDs["tickets"]++
		booking.Tickets[index].BookingID = booking.ID
		booking.Tickets[index].CreatedAt = booking.CreatedAt
		r.db.Tickets[booking.Tickets[index].ID] = booking.Tickets[index]
	}

	r.db.Bookings[booking.ID] = booking
	return booking, nil
}

func (r *BookingRepository) GetByID(id int) (models.Booking, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	booking, ok := r.db.Bookings[id]
	return booking, ok
}

func (r *BookingRepository) GetByUserID(userID int) []models.Booking {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	bookings := make([]models.Booking, 0)
	for _, booking := range r.db.Bookings {
		if booking.UserID == userID {
			bookings = append(bookings, booking)
		}
	}
	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].CreatedAt.After(bookings[j].CreatedAt)
	})
	return bookings
}

func (r *BookingRepository) Cancel(bookingID int) (models.Booking, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	booking, ok := r.db.Bookings[bookingID]
	if !ok {
		return models.Booking{}, errors.New("booking not found")
	}

	if booking.Status == models.BookingPaid {
		return models.Booking{}, errors.New("cannot cancel a paid booking")
	}

	booking.Status = models.BookingCanceled
	r.db.Bookings[bookingID] = booking
	return booking, nil
}

func (r *BookingRepository) MarkPaid(bookingID int) (models.Booking, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	booking, ok := r.db.Bookings[bookingID]
	if !ok {
		return models.Booking{}, errors.New("booking not found")
	}
	if booking.Status == models.BookingPaid {
		return booking, nil
	}

	now := time.Now()
	booking.Status = models.BookingPaid
	booking.PaidAt = &now
	r.db.Bookings[bookingID] = booking
	return booking, nil
}

func (r *BookingRepository) CreatePayment(payment models.Payment) models.Payment {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	payment.ID = r.db.NextIDs["payments"]
	r.db.NextIDs["payments"]++
	payment.PaidAt = time.Now()
	r.db.Payments[payment.ID] = payment
	return payment
}

func (r *BookingRepository) SeatsAvailable(scheduleID int, seatIDs []int) bool {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	return !r.hasConflictingSeatsLocked(scheduleID, seatIDs, 0)
}

func (r *BookingRepository) hasConflictingSeatsLocked(scheduleID int, seatIDs []int, ignoredBookingID int) bool {
	requested := make(map[int]bool, len(seatIDs))
	for _, seatID := range seatIDs {
		requested[seatID] = true
	}

	for _, booking := range r.db.Bookings {
		if booking.ID == ignoredBookingID || booking.ScheduleID != scheduleID || booking.Status == models.BookingCanceled {
			continue
		}
		for _, seatID := range booking.SeatIDs {
			if requested[seatID] {
				return true
			}
		}
	}
	return false
}

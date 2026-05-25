package config

import (
	"fmt"
	"sync"
	"time"

	"gin-M-TIX/models"
)

type Database struct {
	Mu        sync.RWMutex
	Movies    map[int]models.Movie
	Studios   map[int]models.Studio
	Seats     map[int]models.Seat
	Schedules map[int]models.Schedule
	Bookings  map[int]models.Booking
	Tickets   map[int]models.Ticket
	Payments  map[int]models.Payment
	NextIDs   map[string]int
}

func NewDatabase() *Database {
	db := &Database{
		Movies:    make(map[int]models.Movie),
		Studios:   make(map[int]models.Studio),
		Seats:     make(map[int]models.Seat),
		Schedules: make(map[int]models.Schedule),
		Bookings:  make(map[int]models.Booking),
		Tickets:   make(map[int]models.Ticket),
		Payments:  make(map[int]models.Payment),
		NextIDs: map[string]int{
			"movies":    1,
			"studios":   1,
			"seats":     1,
			"schedules": 1,
			"bookings":  1,
			"tickets":   1,
			"payments":  1,
		},
	}

	db.seed()
	return db
}

func (db *Database) seed() {
	db.Movies[1] = models.Movie{ID: 1, Title: "Interstellar", Genre: "Sci-Fi", DurationMinutes: 169, Rating: "PG-13"}
	db.Movies[2] = models.Movie{ID: 2, Title: "The Dark Knight", Genre: "Action", DurationMinutes: 152, Rating: "PG-13"}
	db.NextIDs["movies"] = 3

	db.Studios[1] = models.Studio{ID: 1, Name: "Studio 1", SeatRows: 3, SeatColumns: 5}
	db.Studios[2] = models.Studio{ID: 2, Name: "Studio VIP", SeatRows: 2, SeatColumns: 4}
	db.NextIDs["studios"] = 3

	seatID := 1
	for studioID, studio := range db.Studios {
		for row := 0; row < studio.SeatRows; row++ {
			rowCode := string(rune('A' + row))
			for number := 1; number <= studio.SeatColumns; number++ {
				db.Seats[seatID] = models.Seat{
					ID:       seatID,
					StudioID: studioID,
					Row:      rowCode,
					Number:   number,
					Code:     fmt.Sprintf("%s%d", rowCode, number),
					Status:   models.SeatAvailable,
				}
				seatID++
			}
		}
	}
	db.NextIDs["seats"] = seatID

	now := time.Now()
	db.Schedules[1] = models.Schedule{ID: 1, MovieID: 1, StudioID: 1, StartTime: now.Add(24 * time.Hour), BasePrice: 45000}
	db.Schedules[2] = models.Schedule{ID: 2, MovieID: 2, StudioID: 2, StartTime: now.Add(48 * time.Hour), BasePrice: 75000}
	db.NextIDs["schedules"] = 3
}

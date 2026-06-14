package config

import (
	"fmt"
	"sync"
	"time"

	"gin-M-TIX/models"
)

type Database struct {
	Mu                  sync.RWMutex
	Users               map[int]models.User
	Sessions            map[string]int
	StudentApplications map[int]models.StudentApplication
	Movies              map[int]models.Movie
	Studios             map[int]models.Studio
	Seats               map[int]models.Seat
	Schedules           map[int]models.Schedule
	Bookings            map[int]models.Booking
	Tickets             map[int]models.Ticket
	Payments            map[int]models.Payment
	NextIDs             map[string]int
}

func NewDatabase() *Database {
	db := &Database{}
	db.Reset()
	return db
}

func (db *Database) seed() {
	db.Users[1] = models.User{ID: 1, Username: "admin", Password: "admin", IsAdmin: models.AdminTrue}
	db.Users[2] = models.User{ID: 2, Username: "andi", Password: "andi", IsStudent: models.StudentTrue}
	db.Users[3] = models.User{ID: 3, Username: "budi", Password: "budi"}
	db.Users[4] = models.User{ID: 4, Username: "cici", Password: "cici"}
	db.NextIDs["users"] = 5

	db.Movies[1] = models.Movie{ID: 1, Title: "Interstellar", Genre: "Sci-Fi", DurationMinutes: 169, Rating: "PG-13", PosterURL: "/ui/poster/interstellar.jpeg"}
	db.Movies[2] = models.Movie{ID: 2, Title: "The Dark Knight", Genre: "Action", DurationMinutes: 152, Rating: "PG-13", PosterURL: "/ui/poster/the_dark_knight.jpeg"}
	db.NextIDs["movies"] = 3

	db.Studios[1] = models.Studio{ID: 1, Name: "Studio Candi", SeatRows: 8, SeatColumns: 8}
	db.Studios[2] = models.Studio{ID: 2, Name: "Studio Borobudur", SeatRows: 6, SeatColumns: 6}
	db.NextIDs["studios"] = 3

	seatID := 1
	for studioID := 1; studioID <= len(db.Studios); studioID++ {
		studio := db.Studios[studioID]
		for row := 0; row < studio.SeatRows; row++ {
			for number := 1; number <= studio.SeatColumns; number++ {
				rowCode := string(rune('A' + row))
				db.Seats[seatID] = models.Seat{
					ID:       seatID,
					StudioID: studioID,
					Row:      rowCode,
					Number:   number,
					Code:     fmt.Sprintf("%s%d", rowCode, number),
					IsVIP:    models.SeatVIPStatus(row > studio.SeatRows-3),
					Status:   models.SeatAvailable,
				}
				seatID++
			}
		}
	}
	db.NextIDs["seats"] = seatID

	now := time.Now()
	db.Schedules[1] = models.Schedule{ID: 1, MovieID: 1, StudioID: 1, StartTime: now.Add(96 * time.Hour), BasePrice: 45000}
	db.Schedules[2] = models.Schedule{ID: 2, MovieID: 2, StudioID: 2, StartTime: now.Add(120 * time.Hour), BasePrice: 75000}
	db.NextIDs["schedules"] = 3
}

func (db *Database) Reset() {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	db.Users = make(map[int]models.User)
	db.Sessions = make(map[string]int)
	db.StudentApplications = make(map[int]models.StudentApplication)
	db.Movies = make(map[int]models.Movie)
	db.Studios = make(map[int]models.Studio)
	db.Seats = make(map[int]models.Seat)
	db.Schedules = make(map[int]models.Schedule)
	db.Bookings = make(map[int]models.Booking)
	db.Tickets = make(map[int]models.Ticket)
	db.Payments = make(map[int]models.Payment)
	db.NextIDs = map[string]int{
		"users":     1,
		"movies":    1,
		"studios":   1,
		"seats":     1,
		"schedules": 1,
		"bookings":  1,
		"tickets":   1,
		"payments":  1,
	}

	db.seed()
}

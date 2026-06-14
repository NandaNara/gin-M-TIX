package repositories

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gin-M-TIX/config"
	"gin-M-TIX/models"
)

type StudioRepository struct {
	db *config.Database
}

func NewStudioRepository(db *config.Database) *StudioRepository {
	return &StudioRepository{db: db}
}

func (r *StudioRepository) GetAll() []models.Studio {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	studios := make([]models.Studio, 0, len(r.db.Studios))
	for _, studio := range r.db.Studios {
		studios = append(studios, studio)
	}
	sort.Slice(studios, func(i, j int) bool {
		return studios[i].ID < studios[j].ID
	})
	return studios
}

func (r *StudioRepository) GetByID(id int) (models.Studio, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	studio, ok := r.db.Studios[id]
	return studio, ok
}
func (r *StudioRepository) GetStudios() []models.Studio {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	studios := make([]models.Studio, 0, len(r.db.Studios))
	for _, studio := range r.db.Studios {
		studios = append(studios, studio)
	}
	sort.Slice(studios, func(i, j int) bool {
		return studios[i].ID < studios[j].ID
	})
	return studios
}

func (r *StudioRepository) createStudioSeatsLocked(studio models.Studio) {
	for row := 0; row < studio.SeatRows; row++ {
		rowCode := string(rune('A' + row))
		for number := 1; number <= studio.SeatColumns; number++ {
			seatID := r.db.NextIDs["seats"]
			r.db.NextIDs["seats"]++
			r.db.Seats[seatID] = models.Seat{
				ID:       seatID,
				StudioID: studio.ID,
				Row:      rowCode,
				Number:   number,
				Code:     fmt.Sprintf("%s%d", rowCode, number),
				Status:   models.SeatAvailable,
				IsVIP:    models.SeatVIPStatus(row > studio.SeatRows-3),
			}
		}
	}
}

func (r *StudioRepository) CreateStudio(studio models.Studio) (models.Studio, error) {
	studio.Name = strings.TrimSpace(studio.Name)
	if studio.Name == "" {
		return models.Studio{}, errors.New("studio name is required")
	}
	if studio.SeatRows < 1 || studio.SeatRows > 26 || studio.SeatColumns < 1 || studio.SeatColumns > 50 {
		return models.Studio{}, errors.New("invalid studio data")
	}

	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	studio.ID = r.db.NextIDs["studios"]
	r.db.NextIDs["studios"]++
	r.db.Studios[studio.ID] = studio
	r.createStudioSeatsLocked(studio)
	return studio, nil
}

func (r *StudioRepository) UpdateStudio(id int, studio models.Studio) (models.Studio, error) {
	studio.Name = strings.TrimSpace(studio.Name)
	if studio.Name == "" || studio.SeatRows < 1 || studio.SeatRows > 26 || studio.SeatColumns < 1 || studio.SeatColumns > 50 {
		return models.Studio{}, errors.New("invalid studio data")
	}

	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	current, ok := r.db.Studios[id]
	if !ok {
		return models.Studio{}, errors.New("studio not found")
	}
	for _, schedule := range r.db.Schedules {
		if schedule.StudioID == id && (current.SeatRows != studio.SeatRows || current.SeatColumns != studio.SeatColumns) {
			return models.Studio{}, errors.New("cannot resize a studio used by a schedule")
		}
	}

	studio.ID = id
	r.db.Studios[id] = studio
	if current.SeatRows == studio.SeatRows && current.SeatColumns == studio.SeatColumns {
		return studio, nil
	}

	for seatID, seat := range r.db.Seats {
		if seat.StudioID == id {
			delete(r.db.Seats, seatID)
		}
	}
	r.createStudioSeatsLocked(studio)
	return studio, nil
}

func (r *StudioRepository) DeleteStudio(id int) error {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, ok := r.db.Studios[id]; !ok {
		return errors.New("studio not found")
	}
	for _, schedule := range r.db.Schedules {
		if schedule.StudioID == id {
			return errors.New("studio is used by a schedule")
		}
	}
	for seatID, seat := range r.db.Seats {
		if seat.StudioID == id {
			delete(r.db.Seats, seatID)
		}
	}
	delete(r.db.Studios, id)
	return nil
}

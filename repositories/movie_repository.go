package repositories

import (
	"errors"
	"sort"

	"gin-M-TIX/config"
	"gin-M-TIX/models"
)

type MovieRepository struct {
	db *config.Database
}

func NewMovieRepository(db *config.Database) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetAll() []models.Movie {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	movies := make([]models.Movie, 0, len(r.db.Movies))
	for _, movie := range r.db.Movies {
		movies = append(movies, movie)
	}
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ID < movies[j].ID
	})
	return movies
}

func (r *MovieRepository) GetByID(id int) (models.Movie, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	movie, ok := r.db.Movies[id]
	return movie, ok
}

func (r *MovieRepository) Create(movie models.Movie) models.Movie {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	movie.ID = r.db.NextIDs["movies"]
	r.db.NextIDs["movies"]++
	r.db.Movies[movie.ID] = movie
	return movie
}

func (r *MovieRepository) Update(id int, movie models.Movie) (models.Movie, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, ok := r.db.Movies[id]; !ok {
		return models.Movie{}, errors.New("movie not found")
	}
	movie.ID = id
	r.db.Movies[id] = movie
	return movie, nil
}

func (r *MovieRepository) Delete(id int) error {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, ok := r.db.Movies[id]; !ok {
		return errors.New("movie not found")
	}
	delete(r.db.Movies, id)
	return nil
}

package repositories

import (
	"errors"
	"sort"
	"strings"
	"time"

	"gin-M-TIX/config"
	"gin-M-TIX/models"
)

type UserRepository struct {
	db *config.Database
}

func NewUserRepository(db *config.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, password string) (models.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	username = strings.TrimSpace(username)
	for _, user := range r.db.Users {
		if strings.EqualFold(user.Username, username) {
			return models.User{}, errors.New("username already exists")
		}
	}

	user := models.User{
		ID:       r.db.NextIDs["users"],
		Username: username,
		Password: password,
	}
	r.db.NextIDs["users"]++
	r.db.Users[user.ID] = user
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (models.User, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()
	user, ok := r.db.Users[id]
	return user, ok
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	for _, user := range r.db.Users {
		if strings.EqualFold(user.Username, strings.TrimSpace(username)) {
			return user, true
		}
	}
	return models.User{}, false
}

func (r *UserRepository) SaveSession(token string, userID int) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()
	r.db.Sessions[token] = userID
}

func (r *UserRepository) DeleteSession(token string) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()
	delete(r.db.Sessions, token)
}

func (r *UserRepository) GetUserByToken(token string) (models.User, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	userID, ok := r.db.Sessions[token]
	if !ok {
		return models.User{}, false
	}
	user, ok := r.db.Users[userID]
	return user, ok
}

func (r *UserRepository) SubmitStudentApplication(userID int, filename, contentType, evidencePath string) (models.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	user, ok := r.db.Users[userID]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	if user.IsAdmin == models.AdminTrue {
		return models.User{}, errors.New("admin cannot apply for student status")
	}
	if user.IsStudent == models.StudentTrue {
		return models.User{}, errors.New("user is already a verified student")
	}
	if user.IsStudent == models.StudentPending {
		return models.User{}, errors.New("student application is already pending")
	}

	user.IsStudent = models.StudentPending
	r.db.Users[userID] = user
	r.db.StudentApplications[userID] = models.StudentApplication{
		UserID:       userID,
		Username:     user.Username,
		Filename:     filename,
		ContentType:  contentType,
		EvidencePath: evidencePath,
		SubmittedAt:  time.Now(),
	}
	return user, nil
}

func (r *UserRepository) GetStudentApplication(userID int) (models.StudentApplication, bool) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()
	application, ok := r.db.StudentApplications[userID]
	return application, ok
}

func (r *UserRepository) GetPendingStudentApplications() []models.StudentApplication {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	applications := make([]models.StudentApplication, 0, len(r.db.StudentApplications))
	for _, application := range r.db.StudentApplications {
		applications = append(applications, application)
	}
	sort.Slice(applications, func(i, j int) bool {
		return applications[i].SubmittedAt.Before(applications[j].SubmittedAt)
	})
	return applications
}

func (r *UserRepository) ResolveStudentApplication(userID int, approved bool) (models.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	user, ok := r.db.Users[userID]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	if _, ok := r.db.StudentApplications[userID]; !ok {
		return models.User{}, errors.New("student application not found")
	}

	if approved {
		user.IsStudent = models.StudentTrue
	} else {
		user.IsStudent = models.StudentFalse
	}
	r.db.Users[userID] = user
	delete(r.db.StudentApplications, userID)
	return user, nil
}

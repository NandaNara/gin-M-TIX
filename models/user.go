package models

import "time"

type AdminStatus bool

const (
	AdminTrue  AdminStatus = true
	AdminFalse AdminStatus = false
)

type StudentStatus int

const (
	StudentFalse   StudentStatus = 0
	StudentPending StudentStatus = 1
	StudentTrue    StudentStatus = 2
)

type User struct {
	ID        int           `json:"id"`
	Username  string        `json:"username"`
	Password  string        `json:"-"`
	IsAdmin   AdminStatus   `json:"is_admin"`
	IsStudent StudentStatus `json:"is_student"`
}

type StudentApplication struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Filename     string    `json:"filename"`
	ContentType  string    `json:"content_type"`
	EvidencePath string    `json:"-"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

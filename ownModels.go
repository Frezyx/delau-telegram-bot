package main

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//Task ...
type Task struct {
	ID          int       `json:"id"`
	Time        time.Time `json:"time"`
	Date        time.Time `json:"date"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MarkerIcon  string    `json:"marker_icon"`
	MarkerName  string    `json:"marker_name"`
	IsDone      bool      `json:"is_done"`
	UserID      int       `json:"user_id"`
	Rating      float64   `json:"rating"`
}

//Validate ...
func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.Description, validation.Required),
		validation.Field(&t.MarkerName, validation.Required),
		validation.Field(&t.MarkerIcon, validation.Required),
		validation.Field(&t.Date, validation.Required),
		validation.Field(&t.Time, validation.Required),
	)
}

//AuthRequest ...
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ChatID   int64  `json:"chat_id_tg"`
}

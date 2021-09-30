package book

import (
	"fmt"
	"strings"
	"time"
)

// Book represents book entity
type Book struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int       `json:"rating"`
	Status      Status    `json:"status"`
}

// Status of a book
type Status string

const (
	// StatusCheckedIn represents checked in status
	StatusCheckedIn Status = "checked_in"
	// StatusCheckedOut represents checked out status
	StatusCheckedOut = "checked_out"
)

// ValidationErrors is a collection of book validation errors
type ValidationErrors struct {
	errors []error
}

func (ve *ValidationErrors) Error() string {
	msg := "invalid book entry: "
	var errs []string
	for _, err := range ve.errors {
		errs = append(errs, err.Error())
	}

	return msg + "[" + strings.Join(errs, ", ") + "]"
}

// Validate a book instance for invalid data:
// - invalid rating not in range [1, 3]
// - invalid status
func (b *Book) Validate() error {
	ve := &ValidationErrors{}
	if b.Rating > 3 || b.Rating < 1 {
		err := fmt.Errorf("rating should be between 1 and 3, got %d", b.Rating)
		ve.errors = append(ve.errors, err)
	}
	switch b.Status {
	case StatusCheckedIn:
	case StatusCheckedOut:
		break
	default:
		ve.errors = append(ve.errors, fmt.Errorf("invalid status: %s", b.Status))
	}
	if len(ve.errors) > 0 {
		return ve
	}
	return nil
}

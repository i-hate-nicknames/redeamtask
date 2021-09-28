package book

import (
	"fmt"
	"strings"
	"time"
)

type Book struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int       `json:"rating"`
	Status      Status    `json:"status"`
}

type Status string

const (
	StatusCheckedIn  Status = "checked_in"
	StatusCheckedOut        = "checked_out"
)

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

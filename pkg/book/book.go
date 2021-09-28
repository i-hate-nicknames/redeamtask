package book

import "time"

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

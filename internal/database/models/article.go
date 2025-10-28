package models

import "time"

type Article struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	HTML      string    `db:"html" json:"html"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

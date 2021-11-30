package models

import "time"

type Book struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

package types

import "time"

type Image struct {
	ID        int64      `json:"id"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ImageCreate struct {
	URL string `json:"url"`
}

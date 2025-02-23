package types

import "time"

type Article struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	ShortDescription string     `json:"short_description"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	AuthorID         int64      `json:"author_id"`
	Author           *Author    `json:"author,omitempty"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type ArticleCreate struct {
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	Content          string `json:"content"`
	Status           string `json:"status"`
}

type ArticleUpdate struct {
	Title            string     `json:"title"`
	ShortDescription string     `json:"short_description"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
}

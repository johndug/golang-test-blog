package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Article struct {
	Title            string
	Slug             string
	ShortDescription string
	Content          string
	Status           string
	AuthorID         int64
	PublishedAt      time.Time
}

func Open() (*sql.DB, error) {
	_, err := os.Stat("./skinny_local.db")
	if err == nil {
		os.Remove("./skinny_local.db")
	}
	db, err := sql.Open("sqlite3", "./skinny_local.db")
	if err != nil {
		return nil, err
	}

	if err = createTables(db); err != nil {
		return nil, err
	}

	if err = seedUsers(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			role_id INTEGER NOT NULL,
			last_login DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY (role_id) REFERENCES roles(id)
		)`,
		`CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			bio TEXT,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			short_description TEXT,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'draft',
			author_id INTEGER NOT NULL,
			published_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY (author_id) REFERENCES authors(id)
		)`,
		`CREATE TABLE IF NOT EXISTS article_images (
			article_id INTEGER NOT NULL,
			image_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (article_id, image_id),
			FOREIGN KEY (article_id) REFERENCES articles(id),
			FOREIGN KEY (image_id) REFERENCES images(id)
		)`,
		`CREATE TABLE IF NOT EXISTS images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func seedUsers(db *sql.DB) error {
	// Hash password123 - in production, use bcrypt.GenerateFromPassword
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO roles (name) VALUES (?)`,
		"admin",
	)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO roles (name) VALUES (?)`,
		"author",
	)
	if err != nil {
		return err
	}

	result, err := db.Exec(`
		INSERT INTO users (first_name, last_name, email, password, is_admin, role_id, last_login, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"John",
		"Tavener",
		"john@example.com",
		hashedPassword,
		true,
		2,
		time.Now(),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Printf("Error seeding users: %v", err)
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO authors (first_name, last_name, slug, bio, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"John",
		"Tavener",
		"john-tavener",
		"A sample author bio",
		userID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Printf("Error seeding authors: %v", err)
		return err
	}

	// load multiple articles into the database
	array := []Article{
		{
			Title:            "Article 1",
			Slug:             "article-1",
			ShortDescription: "Short description of article 1",
			Content:          "Content of article 1",
			Status:           "published",
			AuthorID:         userID,
			PublishedAt:      time.Now(),
		},
		{
			Title:            "Article 2",
			Slug:             "article-2",
			ShortDescription: "Short description of article 2",
			Content:          "Content of article 2",
			Status:           "published",
			AuthorID:         userID,
			PublishedAt:      time.Now(),
		},
		{
			Title:            "Article 3",
			Slug:             "article-3",
			ShortDescription: "Short description of article 3",
			Content:          "Content of article 3",
			Status:           "published",
			AuthorID:         userID,
			PublishedAt:      time.Now(),
		},
	}

	for _, article := range array {
		_, err = db.Exec(`
			INSERT INTO articles (title, slug, short_description, content, status, author_id, published_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			article.Title,
			article.Slug,
			article.ShortDescription,
			article.Content,
			article.Status,
			article.AuthorID,
			article.PublishedAt,
		)
		if err != nil {
			log.Printf("Error seeding articles: %v", err)
			return err
		}
	}

	return nil
}

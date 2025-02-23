package stores

import (
	"database/sql"
	"test-ai-api/types"
	"test-ai-api/utils"
	"time"
)

type AuthorStore struct {
	db *sql.DB
}

func NewAuthorStore(db *sql.DB) *AuthorStore {
	return &AuthorStore{db: db}
}

func (s *AuthorStore) Create(author types.AuthorCreate, userID int64) (types.Author, error) {
	slug := utils.GenerateSlug(author.FirstName + " " + author.LastName)
	result, err := s.db.Exec(`
		INSERT INTO authors (first_name, last_name, bio, user_id, slug, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		author.FirstName, author.LastName, author.Bio, userID, slug,
		time.Now(), time.Now(),
	)
	if err != nil {
		return types.Author{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Author{}, err
	}

	return s.GetByID(int64(id))
}

func (s *AuthorStore) GetByID(id int64) (types.Author, error) {
	var author types.Author
	err := s.db.QueryRow(`
		SELECT *
		FROM authors
		WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(
		&author.ID, &author.FirstName, &author.LastName, &author.Slug,
		&author.Bio, &author.UserID, &author.CreatedAt,
		&author.UpdatedAt, &author.DeletedAt,
	)
	if err != nil {
		return types.Author{}, err
	}
	return author, nil
}

func (s *AuthorStore) GetBySlug(slug string) (types.Author, error) {
	var author types.Author
	err := s.db.QueryRow(`
		SELECT *
		FROM authors
		WHERE slug = ? AND deleted_at IS NULL`,
		slug,
	).Scan(
		&author.ID, &author.FirstName, &author.LastName, &author.Slug,
		&author.Bio, &author.UserID, &author.CreatedAt,
		&author.UpdatedAt, &author.DeletedAt,
	)
	if err != nil {
		return types.Author{}, err
	}
	return author, nil
}

func (s *AuthorStore) GetAll(limit int, offset int) ([]types.Author, error) {
	rows, err := s.db.Query(`
		SELECT *
		FROM authors a
		LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []types.Author
	for rows.Next() {
		var author types.Author
		if err := rows.Scan(
			&author.ID, &author.FirstName, &author.LastName, &author.Slug,
			&author.Bio, &author.UserID, &author.CreatedAt,
			&author.UpdatedAt, &author.DeletedAt,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (s *AuthorStore) Update(id int64, author types.AuthorUpdate) (types.Author, error) {
	result, err := s.db.Exec(`
		UPDATE authors SET 
			first_name = ?, last_name = ?, bio = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		author.FirstName, author.LastName, author.Bio, time.Now(), id,
	)
	if err != nil {
		return types.Author{}, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return types.Author{}, err
	}

	return s.GetByID(id)
}

func (s *AuthorStore) Delete(id int64) error {
	_, err := s.db.Exec("UPDATE authors SET deleted_at = ? WHERE id = ?", time.Now(), id)
	return err
}

func (s *AuthorStore) GetByUserID(userID int64) (types.Author, error) {
	var author types.Author
	err := s.db.QueryRow(`
		SELECT id, first_name, last_name, slug, bio, user_id, created_at, updated_at, deleted_at
		FROM authors
		WHERE user_id = ? AND deleted_at IS NULL`,
		userID,
	).Scan(
		&author.ID, &author.FirstName, &author.LastName, &author.Slug,
		&author.Bio, &author.UserID, &author.CreatedAt, &author.UpdatedAt,
		&author.DeletedAt,
	)
	if err != nil {
		return types.Author{}, err
	}
	return author, nil
}

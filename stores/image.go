package stores

import (
	"database/sql"
	"test-ai-api/types"
	"time"
)

type ImageStore struct {
	db *sql.DB
}

func NewImageStore(db *sql.DB) *ImageStore {
	return &ImageStore{db: db}
}

func (s *ImageStore) Create(image types.ImageCreate) (types.Image, error) {
	result, err := s.db.Exec(`
		INSERT INTO images (url, created_at)
		VALUES (?, ?)`,
		image.URL, time.Now(),
	)
	if err != nil {
		return types.Image{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Image{}, err
	}

	return s.GetByID(id)
}

func (s *ImageStore) GetByID(id int64) (types.Image, error) {
	var image types.Image
	err := s.db.QueryRow(`
		SELECT id, url, created_at, deleted_at
		FROM images
		WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(&image.ID, &image.URL, &image.CreatedAt, &image.DeletedAt)
	if err != nil {
		return types.Image{}, err
	}
	return image, nil
}

func (s *ImageStore) Delete(id int64) error {
	_, err := s.db.Exec(
		"UPDATE images SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL",
		time.Now(), id,
	)
	return err
}

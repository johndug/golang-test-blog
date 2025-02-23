package stores

import (
	"database/sql"
	"test-ai-api/types"
	"test-ai-api/utils"
	"time"
)

type ArticleStore struct {
	db *sql.DB
}

func NewArticleStore(db *sql.DB) *ArticleStore {
	return &ArticleStore{db: db}
}

func (s *ArticleStore) Create(article types.ArticleCreate, authorID int64) (types.Article, error) {
	slug := utils.GenerateSlug(article.Title)
	result, err := s.db.Exec(`
		INSERT INTO articles (title, slug, short_description, content, status, author_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		article.Title, slug, article.ShortDescription, article.Content,
		article.Status, authorID, time.Now(), time.Now(),
	)
	if err != nil {
		return types.Article{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Article{}, err
	}

	return s.GetByID(id)
}

func (s *ArticleStore) GetByID(id int64) (types.Article, error) {
	var article types.Article
	var author types.Author
	err := s.db.QueryRow(`
		SELECT a.*, au.*
		FROM articles a
		LEFT JOIN authors au ON a.author_id = au.id
		WHERE a.id = ? AND a.deleted_at IS NULL`,
		id,
	).Scan(
		&article.ID, &article.Title, &article.Slug, &article.ShortDescription,
		&article.Content, &article.Status, &article.AuthorID, &article.PublishedAt,
		&article.CreatedAt, &article.UpdatedAt, &article.DeletedAt,
		&author.ID, &author.FirstName, &author.LastName,
		&author.Slug, &author.Bio, &author.UserID,
		&author.CreatedAt, &author.UpdatedAt, &author.DeletedAt,
	)
	if err != nil {
		return types.Article{}, err
	}
	article.Author = &author
	return article, nil
}

func (s *ArticleStore) GetAll(limit int, offset int) ([]types.Article, error) {
	rows, err := s.db.Query(`
		SELECT a.*, 
			   au.*
		FROM articles a
		LEFT JOIN authors au ON a.author_id = au.id
		WHERE a.deleted_at IS NULL
		LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return []types.Article{}, err
	}
	defer rows.Close()

	var articles []types.Article

	for rows.Next() {
		var article types.Article
		var author types.Author
		if err := rows.Scan(
			&article.ID, &article.Title, &article.Slug, &article.ShortDescription,
			&article.Content, &article.Status, &article.AuthorID,
			&article.PublishedAt, &article.CreatedAt, &article.UpdatedAt,
			&article.DeletedAt, &author.ID, &author.FirstName, &author.LastName,
			&author.Slug, &author.Bio, &author.UserID,
			&author.CreatedAt, &author.UpdatedAt, &author.DeletedAt,
		); err != nil {
			return []types.Article{}, err
		}
		article.Author = &author
		articles = append(articles, article)
	}
	return articles, nil
}

func (s *ArticleStore) Update(id int64, article types.ArticleUpdate) (types.Article, error) {
	result, err := s.db.Exec(`
		UPDATE articles SET 
			title = ?, short_description = ?, content = ?, status = ?, published_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		article.Title, article.ShortDescription, article.Content,
		article.Status, article.PublishedAt, time.Now(), id,
	)
	if err != nil {
		return types.Article{}, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return types.Article{}, err
	}

	return s.GetByID(id)
}

func (s *ArticleStore) Delete(id int64) error {
	_, err := s.db.Exec("UPDATE articles SET deleted_at = ? WHERE id = ?", time.Now(), id)
	return err
}

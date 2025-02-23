package stores

import (
	"database/sql"
	"fmt"
	"test-ai-api/types"

	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetAll(limit int, offset int) ([]types.User, error) {
	users, err := s.db.Query("SELECT * FROM users WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer users.Close()

	usersList := []types.User{}
	for users.Next() {
		var user types.User
		err := users.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func (s *UserStore) GetByEmail(email string) (types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return types.User{}, err
	}
	defer rows.Close()

	var user types.User
	err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (s *UserStore) GetByID(id int64) (types.User, error) {
	var user types.User
	var role types.Role
	err := s.db.QueryRow(`
		SELECT u.*, r.id, r.name 
		FROM users u 
		INNER JOIN roles r ON u.role_id = r.id 
		WHERE u.id = ? AND u.deleted_at IS NULL`,
		id,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.IsAdmin, &user.RoleID, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		&user.DeletedAt, &role.ID, &role.Name)
	if err != nil {
		return types.User{}, err
	}
	user.Role = role
	return user, nil
}

func (s *UserStore) Create(user types.User) (types.User, error) {
	rows, err := s.db.Query("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return types.User{}, err
	}
	defer rows.Close()

	return user, nil
}

func (s *UserStore) Update(user types.User) (types.User, error) {
	rows, err := s.db.Query("UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?", user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return types.User{}, err
	}
	defer rows.Close()

	return user, nil
}

func (s *UserStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Register(ur types.UserRegister) (types.User, error) {
	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ur.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, err
	}

	result, err := s.db.Exec(
		"INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
		ur.FirstName, ur.LastName, ur.Email, string(hashedPassword),
	)
	if err != nil {
		return types.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.User{}, err
	}

	user, err := s.GetByID(id)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (s *UserStore) Login(email, password string) (types.User, error) {
	var user types.User
	var role types.Role
	var hashedPassword string

	err := s.db.QueryRow(
		`SELECT users.id, users.first_name, users.last_name, users.email, users.password, roles.id, roles.name
		FROM users 
		INNER JOIN roles ON users.role_id = roles.id
		WHERE users.email = ? AND users.deleted_at IS NULL`,
		email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &hashedPassword, &role.ID, &role.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("invalid credentials")
		}
		return types.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return types.User{}, fmt.Errorf("invalid credentials")
	}
	user.Role = role

	return user, nil
}

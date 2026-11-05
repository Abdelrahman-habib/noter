package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) (int, error)
	Authenticate(email, password string) (int, error)
	GetByID(id int) (User, error)
	Exists(id int) (bool, error)
	ChangePassword(id int, currentPassword, newPassword string) error
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, CURRENT_TIMESTAMP) RETURNING id`

	var id int
	err = m.DB.QueryRow(stmt, name, email, string(hashedPassword)).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgUniqueViolationCode && pgErr.ConstraintName == "users_uc_email" {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte

	stmt := `SELECT id, hashed_password FROM users WHERE lower(email) = lower($1)`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashed_password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) GetByID(id int) (User, error) {
	stmt := `SELECT id, name, email, created FROM users WHERE id = $1`
	var user User
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte
	stmt := `SELECT hashed_password FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newPassword), 12)

	if err != nil {
		return err
	}

	stmt = `UPDATE users SET hashed_password = $1 WHERE id = $2`

	_, err = m.DB.Exec(stmt, string(hashedPassword), id)

	return err
}

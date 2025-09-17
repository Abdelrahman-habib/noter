package mocks

import (
	"time"

	"github.com/Abdelrahman-habib/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	switch email {
	case "dupe@example.com":
		return 0, models.ErrDuplicateEmail
	default:
		return 2, nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) GetByID(id int) (models.User, error) {
	switch id {
	case 1:
		return models.User{
			Email:   "alice@example.com",
			ID:      1,
			Name:    "alice",
			Created: time.Date(2012, 2, 2, 12, 10, 0, 0, time.Local),
		}, nil
	default:
		return models.User{}, nil
	}
}

func (m *UserModel) ChangePassword(id int, current, new string) error {
	switch id {
	case 1:
		if current != "pa$$word" {
			return models.ErrInvalidCredentials
		}
		return nil
	default:
		return models.ErrNoRecord
	}
}

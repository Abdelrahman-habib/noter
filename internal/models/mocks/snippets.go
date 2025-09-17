package mocks

import (
	"time"

	"github.com/Abdelrahman-habib/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:        "550e8400-e29b-41d4-a716-446655440000",
	Title:     "An old silent pond",
	Content:   "An old silent pond...",
	Created:   time.Now(),
	Expires:   time.Now(),
	Public:    true,
	CreatedBy: 1,
}

var mockSnippetWithUsername = models.SnippetWithUsername{
	Snippet:  mockSnippet,
	Username: "John Doe",
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int, public bool, createdBy int) (string, error) {
	return "550e8400-e29b-41d4-a716-446655440001", nil
}
func (m *SnippetModel) Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error) {
	return "550e8400-e29b-41d4-a716-446655440001", nil
}
func (m *SnippetModel) Get(id string, createdBy *int) (models.SnippetWithUsername, error) {
	switch id {
	case "550e8400-e29b-41d4-a716-446655440000":
		return mockSnippetWithUsername, nil
	default:
		return models.SnippetWithUsername{}, models.ErrNoRecord
	}
}
func (m *SnippetModel) Latest() ([]models.SnippetWithUsername, error) {
	return []models.SnippetWithUsername{mockSnippetWithUsername}, nil
}

func (m *SnippetModel) GetByPage(page int, limit int, public *bool, createdBy *int) ([]models.SnippetWithUsername, models.PaginationMetaData, error) {
	return []models.SnippetWithUsername{mockSnippetWithUsername}, models.PaginationMetaData{
		HasNext: false,
	}, nil
}

func (m *SnippetModel) GetTotalPages(public *bool, createdBy *int) (int, error) {
	return 1, nil
}

func (m *SnippetModel) Delete(id string, createdBy *int) error {
	return nil
}

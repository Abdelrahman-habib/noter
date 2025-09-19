package mocks

import (
	"time"

	"github.com/Abdelrahman-habib/noter/internal/models"
)

var mockNote = models.Note{
	ID:        "550e8400-e29b-41d4-a716-446655440000",
	Title:     "An old silent pond",
	Content:   "An old silent pond...",
	Created:   time.Now(),
	Expires:   time.Now(),
	Public:    true,
	CreatedBy: 1,
}

var mockNoteWithUsername = models.NoteWithUsername{
	Note:     mockNote,
	Username: "John Doe",
}

type NoteModel struct{}

func (m *NoteModel) Insert(title string, content string, expires int, public bool, createdBy int) (string, error) {
	return "550e8400-e29b-41d4-a716-446655440001", nil
}
func (m *NoteModel) Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error) {
	return "550e8400-e29b-41d4-a716-446655440001", nil
}
func (m *NoteModel) Get(id string, createdBy *int) (models.NoteWithUsername, error) {
	switch id {
	case "550e8400-e29b-41d4-a716-446655440000":
		return mockNoteWithUsername, nil
	default:
		return models.NoteWithUsername{}, models.ErrNoRecord
	}
}
func (m *NoteModel) Latest() ([]models.NoteWithUsername, error) {
	return []models.NoteWithUsername{mockNoteWithUsername}, nil
}

func (m *NoteModel) GetByPage(page int, limit int, public *bool, createdBy *int) ([]models.NoteWithUsername, models.PaginationMetaData, error) {
	return []models.NoteWithUsername{mockNoteWithUsername}, models.PaginationMetaData{
		HasNext: false,
	}, nil
}

func (m *NoteModel) GetTotalPages(public *bool, createdBy *int) (int, error) {
	return 1, nil
}

func (m *NoteModel) Delete(id string, createdBy *int) error {
	return nil
}

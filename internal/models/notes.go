package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Define a Note type to hold the data for an individual note. Notice how
// the fields of the struct correspond to the fields in our MySQL notes
// table?
type Note struct {
	ID        string
	Title     string
	Content   string
	Created   time.Time
	Public    bool
	CreatedBy int
	Expires   time.Time
}

type NoteWithUsername struct {
	Note
	Username string
}

type NotesFilters struct {
	ShowPublic *bool
}

// Define a NoteModel type which wraps a sql.DB connection pool.
type NoteModel struct {
	DB *sql.DB
}

type PaginationMetaData struct {
	HasNext bool
}

type NoteModelInterface interface {
	Insert(title string, content string, expires int, public bool, createdBy int) (string, error)
	Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error)
	Get(id string, createdBy *int) (NoteWithUsername, error)
	Latest() ([]NoteWithUsername, error)
	GetByPage(page int, limit int, public *bool, createdBy *int) ([]NoteWithUsername, PaginationMetaData, error)
	GetTotalPages(public *bool, createdBy *int) (int, error)
	Delete(id string, createdBy *int) error
}

// This will insert a new notes into the database.

func (m *NoteModel) Insert(title string, content string, expires int, public bool, createdBy int) (string, error) {
	stmt := `INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?)`
	id := uuid.New().String()

	_, err := m.DB.Exec(stmt, id, title, content, expires, public, createdBy)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *NoteModel) Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error) {
	stmt := `UPDATE notes SET title = ?, content = ?, expires = DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), public = ?, created_by = ? WHERE id = ? AND created_by = ?`
	_, err := m.DB.Exec(stmt, title, content, expires, public, createdBy, id, createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		}
		return "", err
	}
	return id, nil
}

// This will return a specific note based on its id.
func (m *NoteModel) Get(id string, createdBy *int) (NoteWithUsername, error) {
	stmt := `SELECT notes.id, notes.title, notes.content, notes.created, notes.expires, notes.public, notes.created_by, users.name as username 
	FROM notes 
	JOIN users ON notes.created_by = users.id 
	WHERE notes.expires > UTC_TIMESTAMP() AND notes.id = ?`

	var args []interface{}
	args = append(args, id)

	// Access control: if no user ID provided, only show public notes
	// if user ID provided, show public notes OR private notes created by that user
	if createdBy != nil {
		stmt += ` AND (notes.public = TRUE OR notes.created_by = ?)`
		args = append(args, *createdBy)
	} else {
		stmt += ` AND notes.public = TRUE`
	}

	var s NoteWithUsername

	err := m.DB.QueryRow(stmt, args...).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NoteWithUsername{}, ErrNoRecord
		}
		return NoteWithUsername{}, err
	}
	return s, nil
}

// This will return the 10 most recently created public notes.
func (m *NoteModel) Latest() ([]NoteWithUsername, error) {
	stmt := `SELECT notes.id, notes.title, notes.content, notes.created, notes.expires, notes.public, notes.created_by, users.name as username 
	FROM notes
	JOIN users ON notes.created_by = users.id
	WHERE notes.expires > UTC_TIMESTAMP() AND notes.public = TRUE ORDER BY notes.id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []NoteWithUsername
	for rows.Next() {
		var s NoteWithUsername
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
		if err != nil {
			return nil, err
		}
		notes = append(notes, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (m *NoteModel) GetTotalPages(public *bool, createdBy *int) (int, error) {
	stmt := `SELECT COUNT(*) FROM notes WHERE expires > UTC_TIMESTAMP()`

	// Build dynamic query parameters
	var args []interface{}

	// Filter by created_by if specified (shows only notes by that user)
	if createdBy != nil {
		stmt += ` AND notes.created_by = ?`
		args = append(args, *createdBy)
	}

	// Filter by public/private if specified
	if public != nil {
		stmt += ` AND notes.public = ?`
		args = append(args, *public)
	}

	var total int
	err := m.DB.QueryRow(stmt, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (m *NoteModel) GetByPage(page int, limit int, public *bool, createdBy *int) ([]NoteWithUsername, PaginationMetaData, error) {
	originalLimit := limit
	limit = limit + 1 // to check if there is a next page
	offset := (page - 1) * originalLimit
	stmt := `SELECT notes.id, notes.title, notes.content, notes.created, notes.expires, notes.public, notes.created_by, users.name as username 
	FROM notes
	JOIN users ON notes.created_by = users.id
	WHERE notes.expires > UTC_TIMESTAMP()`

	// Build dynamic query parameters
	var args []interface{}

	// Filter by created_by if specified (shows only notes by that user)
	if createdBy != nil {
		stmt += ` AND notes.created_by = ?`
		args = append(args, *createdBy)
	}

	// Filter by public/private if specified
	if public != nil {
		stmt += ` AND notes.public = ?`
		args = append(args, *public)
	}

	stmt += ` ORDER BY notes.id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	meta := PaginationMetaData{
		HasNext: false,
	}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	var notes []NoteWithUsername
	for rows.Next() {
		var s NoteWithUsername
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
		if err != nil {
			return nil, meta, err
		}
		notes = append(notes, s)
	}

	if len(notes) > originalLimit {
		meta.HasNext = true
		notes = notes[:originalLimit]
	}

	return notes, meta, nil
}

func (m *NoteModel) Delete(id string, createdBy *int) error {
	stmt := `DELETE FROM notes WHERE id = ? AND created_by = ?`
	_, err := m.DB.Exec(stmt, id, *createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	return nil
}

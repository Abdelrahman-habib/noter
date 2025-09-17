package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID        string
	Title     string
	Content   string
	Created   time.Time
	Public    bool
	CreatedBy int
	Expires   time.Time
}

type SnippetWithUsername struct {
	Snippet
	Username string
}

type SnippetsFilters struct {
	ShowPublic *bool
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

type PaginationMetaData struct {
	HasNext bool
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int, public bool, createdBy int) (string, error)
	Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error)
	Get(id string, createdBy *int) (SnippetWithUsername, error)
	Latest() ([]SnippetWithUsername, error)
	GetByPage(page int, limit int, public *bool, createdBy *int) ([]SnippetWithUsername, PaginationMetaData, error)
	GetTotalPages(public *bool, createdBy *int) (int, error)
	Delete(id string, createdBy *int) error
}

// This will insert a new snippet into the database.

func (m *SnippetModel) Insert(title string, content string, expires int, public bool, createdBy int) (string, error) {
	stmt := `INSERT INTO snippets (id, title, content, created, expires, public, created_by) VALUES (?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?)`
	id := uuid.New().String()

	_, err := m.DB.Exec(stmt, id, title, content, expires, public, createdBy)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *SnippetModel) Update(id string, title string, content string, expires int, public bool, createdBy int) (string, error) {
	stmt := `UPDATE snippets SET title = ?, content = ?, expires = DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), public = ?, created_by = ? WHERE id = ? AND created_by = ?`
	_, err := m.DB.Exec(stmt, title, content, expires, public, createdBy, id, createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		}
		return "", err
	}
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id string, createdBy *int) (SnippetWithUsername, error) {
	stmt := `SELECT snippets.id, snippets.title, snippets.content, snippets.created, snippets.expires, snippets.public, snippets.created_by, users.name as username 
	FROM snippets 
	JOIN users ON snippets.created_by = users.id 
	WHERE snippets.expires > UTC_TIMESTAMP() AND snippets.id = ?`

	var args []interface{}
	args = append(args, id)

	// Access control: if no user ID provided, only show public snippets
	// if user ID provided, show public snippets OR private snippets created by that user
	if createdBy != nil {
		stmt += ` AND (snippets.public = TRUE OR snippets.created_by = ?)`
		args = append(args, *createdBy)
	} else {
		stmt += ` AND snippets.public = TRUE`
	}

	var s SnippetWithUsername

	err := m.DB.QueryRow(stmt, args...).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SnippetWithUsername{}, ErrNoRecord
		}
		return SnippetWithUsername{}, err
	}
	return s, nil
}

// This will return the 10 most recently created public snippets.
func (m *SnippetModel) Latest() ([]SnippetWithUsername, error) {
	stmt := `SELECT snippets.id, snippets.title, snippets.content, snippets.created, snippets.expires, snippets.public, snippets.created_by, users.name as username 
	FROM snippets
	JOIN users ON snippets.created_by = users.id
	WHERE snippets.expires > UTC_TIMESTAMP() AND snippets.public = TRUE ORDER BY snippets.id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []SnippetWithUsername
	for rows.Next() {
		var s SnippetWithUsername
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) GetTotalPages(public *bool, createdBy *int) (int, error) {
	stmt := `SELECT COUNT(*) FROM snippets WHERE expires > UTC_TIMESTAMP()`

	// Build dynamic query parameters
	var args []interface{}

	// Filter by created_by if specified (shows only snippets by that user)
	if createdBy != nil {
		stmt += ` AND snippets.created_by = ?`
		args = append(args, *createdBy)
	}

	// Filter by public/private if specified
	if public != nil {
		stmt += ` AND snippets.public = ?`
		args = append(args, *public)
	}

	var total int
	err := m.DB.QueryRow(stmt, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (m *SnippetModel) GetByPage(page int, limit int, public *bool, createdBy *int) ([]SnippetWithUsername, PaginationMetaData, error) {
	originalLimit := limit
	limit = limit + 1 // to check if there is a next page
	offset := (page - 1) * originalLimit
	stmt := `SELECT snippets.id, snippets.title, snippets.content, snippets.created, snippets.expires, snippets.public, snippets.created_by, users.name as username 
	FROM snippets
	JOIN users ON snippets.created_by = users.id
	WHERE snippets.expires > UTC_TIMESTAMP()`

	// Build dynamic query parameters
	var args []interface{}

	// Filter by created_by if specified (shows only snippets by that user)
	if createdBy != nil {
		stmt += ` AND snippets.created_by = ?`
		args = append(args, *createdBy)
	}

	// Filter by public/private if specified
	if public != nil {
		stmt += ` AND snippets.public = ?`
		args = append(args, *public)
	}

	stmt += ` ORDER BY snippets.id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	meta := PaginationMetaData{
		HasNext: false,
	}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	var snippets []SnippetWithUsername
	for rows.Next() {
		var s SnippetWithUsername
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public, &s.CreatedBy, &s.Username)
		if err != nil {
			return nil, meta, err
		}
		snippets = append(snippets, s)
	}

	if len(snippets) > originalLimit {
		meta.HasNext = true
		snippets = snippets[:originalLimit]
	}

	return snippets, meta, nil
}

func (m *SnippetModel) Delete(id string, createdBy *int) error {
	stmt := `DELETE FROM snippets WHERE id = ? AND created_by = ?`
	_, err := m.DB.Exec(stmt, id, *createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	return nil
}

package models

import (
	"errors"
	"testing"
	"time"

	"github.com/Abdelrahman-habib/noter/internal/assert"
)

// Seed notes: two public, one private, all owned by user 1.
var seedNoteIDs = []string{
	"550e8400-e29b-41d4-a716-446655440000",
	"550e8400-e29b-41d4-a716-446655440001",
	"550e8400-e29b-41d4-a716-446655440002",
}

func TestNoteModelInsertAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	t.Run("Insert returns generated UUID id", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		id, err := m.Insert("Integration note", "Some content", 10, true, 1)
		assert.NilError(t, err)
		assert.Equal(t, len(id) == 36, true)

		note, err := m.Get(id, intPtr(1))
		assert.NilError(t, err)
		assert.Equal(t, note.ID, id)
		assert.Equal(t, note.Title, "Integration note")
		assert.Equal(t, note.Content, "Some content")
		assert.Equal(t, note.Public, true)
		assert.Equal(t, note.CreatedBy, 1)
		assert.Equal(t, note.Username, "Alice Jones")

		assert.Equal(t, note.Expires.Sub(note.Created), 240*time.Hour)
	})

	t.Run("Unauthenticated sees only public notes", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		_, err := m.Get(seedNoteIDs[2], nil)
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)

		note, err := m.Get(seedNoteIDs[0], nil)
		assert.NilError(t, err)
		assert.Equal(t, note.ID, seedNoteIDs[0])
	})

	t.Run("Expired note with zero days is not visible", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		id, err := m.Insert("Expiring now", "Content", 0, true, 1)
		assert.NilError(t, err)

		_, err = m.Get(id, intPtr(1))
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)
	})

	t.Run("Private note visible to its creator only", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		_, err := m.Get(seedNoteIDs[2], intPtr(2))
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)

		note, err := m.Get(seedNoteIDs[2], intPtr(1))
		assert.NilError(t, err)
		assert.Equal(t, note.CreatedBy, 1)
		assert.Equal(t, note.Expires.Sub(note.Created), 7*24*time.Hour)
	})

	t.Run("Unknown id returns ErrNoRecord", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		_, err := m.Get("0197dead-beef-41d4-a716-44665544dead", nil)
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)
	})
}

func TestNoteModelLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	t.Run("Excludes expired public notes", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		_, err := m.Insert("[expired] public note", "invisible", 0, true, 1)
		assert.NilError(t, err)

		notes, err := m.Latest()
		assert.NilError(t, err)
		assert.Equal(t, len(notes), 2)
		for _, note := range notes {
			assert.Equal(t, note.Public, true)
		}
	})
}

func TestNoteModelUpdateAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	t.Run("Creator can update and delete", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		var before, after time.Time
		assert.NilError(t, m.DB.QueryRow(`SELECT current_timestamp`).Scan(&before))

		id, err := m.Update(seedNoteIDs[0], "Updated title", "Updated content", 30, true, 1)
		assert.NilError(t, err)
		assert.Equal(t, id, seedNoteIDs[0])

		assert.NilError(t, m.DB.QueryRow(`SELECT current_timestamp`).Scan(&after))

		note, err := m.Get(id, nil)
		assert.NilError(t, err)
		assert.Equal(t, note.Title, "Updated title")
		assert.Equal(t, note.Content, "Updated content")

		// Expiry is measured from the update time, not the creation time:
		// every timestamp comes from the same server clock, so a creation-based
		// expiry would land seconds before the lower bound.
		assert.Equal(t, !note.Expires.Before(before.Add(30*24*time.Hour)) && !note.Expires.After(after.Add(30*24*time.Hour)), true)

		err = m.Delete(id, intPtr(1))
		assert.NilError(t, err)

		_, err = m.Get(id, nil)
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)
	})

	t.Run("Non-creator cannot update", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		// Non-owner updates match no rows: the model reports ErrNoRecord so
		// handlers can respond with 404.
		_, err := m.Update(seedNoteIDs[0], "Hijacked title", "Hijacked content", 30, true, 999)
		assert.Equal(t, errors.Is(err, ErrNoRecord), true)
	})
}

func TestNoteModelFiltersAndPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	// Add a fourth note to exercise pagination.
	addNote := func(t *testing.T, m *NoteModel) string {
		t.Helper()
		id, err := m.Insert("[extra] Pagination note", "More content", 30, true, 1)
		assert.NilError(t, err)
		return id
	}

	tests := []struct {
		name          string
		public        *bool
		createdBy     *int
		wantTotal     int
		wantPage1     int
		wantPage1Next bool
		wantPage2     int
		wantPage2Next bool
	}{
		{
			name:          "no filters",
			public:        nil,
			createdBy:     nil,
			wantTotal:     4,
			wantPage1:     2,
			wantPage1Next: true,
			wantPage2:     2,
			wantPage2Next: false,
		},
		{
			name:          "public only",
			public:        boolPtr(true),
			createdBy:     nil,
			wantTotal:     3,
			wantPage1:     2,
			wantPage1Next: true,
			wantPage2:     1,
			wantPage2Next: false,
		},
		{
			name:          "created by only",
			public:        nil,
			createdBy:     intPtr(1),
			wantTotal:     4,
			wantPage1:     2,
			wantPage1Next: true,
			wantPage2:     2,
			wantPage2Next: false,
		},
		{
			name:          "private notes created by user 1",
			public:        boolPtr(false),
			createdBy:     intPtr(1),
			wantTotal:     1,
			wantPage1:     1,
			wantPage1Next: false,
			wantPage2:     0,
			wantPage2Next: false,
		},
		{
			name:          "private notes created by other user",
			public:        boolPtr(false),
			createdBy:     intPtr(2),
			wantTotal:     0,
			wantPage1:     0,
			wantPage1Next: false,
			wantPage2:     0,
			wantPage2Next: false,
		},
		{
			name:          "public notes created by user 1",
			public:        boolPtr(true),
			createdBy:     intPtr(1),
			wantTotal:     3,
			wantPage1:     2,
			wantPage1Next: true,
			wantPage2:     1,
			wantPage2Next: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := NoteModel{db}
			addNote(t, &m)

			total, err := m.GetTotalPages(tt.public, tt.createdBy)
			assert.NilError(t, err)
			assert.Equal(t, total, tt.wantTotal)

			page1, meta, err := m.GetByPage(1, 2, tt.public, tt.createdBy)
			assert.NilError(t, err)
			assert.Equal(t, len(page1), tt.wantPage1)
			assert.Equal(t, meta.HasNext, tt.wantPage1Next)

			page2, meta, err := m.GetByPage(2, 2, tt.public, tt.createdBy)
			assert.NilError(t, err)
			assert.Equal(t, len(page2), tt.wantPage2)
			assert.Equal(t, meta.HasNext, tt.wantPage2Next)
		})
	}

	t.Run("HasNext is reported when more rows exist", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}
		id := addNote(t, &m)

		notes, meta, err := m.GetByPage(1, 1, nil, nil)
		assert.NilError(t, err)
		assert.Equal(t, len(notes), 1)
		// Newest-created note first: ordering is by created, not the random UUID id.
		assert.Equal(t, notes[0].ID, id)
		assert.Equal(t, meta.HasNext, true)

		notes, meta, err = m.GetByPage(2, 1, nil, nil)
		assert.NilError(t, err)
		assert.Equal(t, len(notes), 1)
		assert.Equal(t, meta.HasNext, true)

		notes, meta, err = m.GetByPage(4, 1, nil, nil)
		assert.NilError(t, err)
		assert.Equal(t, len(notes), 1)
		assert.Equal(t, meta.HasNext, false)
	})

	t.Run("GetTotalPages does not count expired notes", func(t *testing.T) {
		db := newTestDB(t)
		m := NoteModel{db}

		_, err := m.Insert("[expired] note", "invisible", 0, true, 1)
		assert.NilError(t, err)

		total, err := m.GetTotalPages(nil, nil)
		assert.NilError(t, err)
		assert.Equal(t, total, 3)
	})
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

package models

import (
	"errors"
	"testing"

	"github.com/Abdelrahman-habib/noter/internal/assert"
	"github.com/pressly/goose/v3"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test.
			db := newTestDB(t)
			// Create a new instance of the UserModel.
			m := UserModel{db}
			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	t.Run("Signup returns generated id", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		id, err := m.Insert("Test Person", "test-person@example.com", "emails-expire-soon!")
		assert.NilError(t, err)
		assert.Equal(t, id != 0, true)

		// The seed user takes ID 1.
		assert.Equal(t, id >= 2, true)

		exists, err := m.Exists(id)
		assert.NilError(t, err)
		assert.Equal(t, exists, true)
	})

	t.Run("Duplicate email is rejected", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		// Same email as the seed user, different case.
		id, err := m.Insert("Alice Copy", "ALICE@EXAMPLE.COM", "emails-expire-soon!")
		assert.Equal(t, id, 0)
		assert.Equal(t, errors.Is(err, ErrDuplicateEmail), true)
	})
}

func TestUserModelAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	t.Run("Valid credentials with mixed-case email", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		id, err := m.Insert("Case Person", "Case.Person@Example.COM", "emails-expire-soon!")
		assert.NilError(t, err)

		for _, email := range []string{
			"Case.Person@Example.COM",
			"case.person@example.com",
			"CASE.PERSON@EXAMPLE.COM",
		} {
			got, err := m.Authenticate(email, "emails-expire-soon!")
			assert.NilError(t, err)
			assert.Equal(t, got, id)
		}
	})

	t.Run("Wrong password returns ErrInvalidCredentials", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		_, err := m.Authenticate("alice@example.com", "wrong-password")
		assert.Equal(t, errors.Is(err, ErrInvalidCredentials), true)
	})

	t.Run("Unknown email returns ErrInvalidCredentials", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		_, err := m.Authenticate("nobody@example.com", "emails-expire-soon!")
		assert.Equal(t, errors.Is(err, ErrInvalidCredentials), true)
	})
}

func TestUserModelChangePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := UserModel{db}

	id, err := m.Insert("Change Person", "change-person@example.com", "first-password!")
	assert.NilError(t, err)

	user, err := m.GetByID(id)
	assert.NilError(t, err)
	assert.Equal(t, user.ID, id)
	assert.Equal(t, user.Name, "Change Person")
	assert.Equal(t, user.Email, "change-person@example.com")
	assert.Equal(t, user.Created.IsZero(), false)

	err = m.ChangePassword(id, "wrong-current-password", "second-password!")
	assert.Equal(t, errors.Is(err, ErrInvalidCredentials), true)

	err = m.ChangePassword(id, "first-password!", "second-password!")
	assert.NilError(t, err)

	got, err := m.Authenticate("change-person@example.com", "second-password!")
	assert.NilError(t, err)
	assert.Equal(t, got, id)

	_, err = m.Authenticate("change-person@example.com", "first-password!")
	assert.Equal(t, errors.Is(err, ErrInvalidCredentials), true)
}

// Reseeding must preserve existing users and assign new IDs.
func TestUserModelSeedReseed(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := UserModel{db}

	exists, err := m.Exists(1)
	assert.NilError(t, err)
	assert.Equal(t, exists, true)

	secondID, err := m.Insert("Second Person", "second-person@example.com", "emails-expire-soon!")
	assert.NilError(t, err)

	if err := goose.Reset(db, "seed", goose.WithNoVersioning()); err != nil {
		t.Fatalf("failed to reset seeds: %v", err)
	}
	if err := goose.Up(db, "seed", goose.WithNoVersioning()); err != nil {
		t.Fatalf("failed to re-run seeds: %v", err)
	}

	var seedID int
	err = db.QueryRow(`SELECT id FROM users WHERE email = 'alice@example.com'`).Scan(&seedID)
	assert.NilError(t, err)
	assert.Equal(t, seedID > secondID, true)

	exists, err = m.Exists(1)
	assert.NilError(t, err)
	assert.Equal(t, exists, false)

	notes := NoteModel{db}

	// Reseeded notes belong to the new seed user.
	note, err := notes.Get(seedNoteIDs[2], intPtr(seedID))
	assert.NilError(t, err)
	assert.Equal(t, note.CreatedBy, seedID)

	_, err = notes.Get(seedNoteIDs[2], intPtr(1))
	assert.Equal(t, errors.Is(err, ErrNoRecord), true)

	thirdID, err := m.Insert("Third Person", "third-person@example.com", "emails-expire-soon!")
	assert.NilError(t, err)
	assert.Equal(t, thirdID != seedID && thirdID != secondID, true)
}

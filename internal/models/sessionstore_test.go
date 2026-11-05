package models

import (
	"bytes"
	"testing"
	"time"

	"github.com/Abdelrahman-habib/noter/internal/assert"
	"github.com/alexedwards/scs/postgresstore"
)

// Check the session store against our migrated schema.
func TestPostgresSessionStoreRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)

	store := postgresstore.New(db)
	defer store.StopCleanup()

	token := "sctest-token-0123456789abcdefghijklmnopqrstuvwxyz"
	expiry := time.Now().Add(12 * time.Hour)

	payload := []byte("roundtrip session payload")
	err := store.Commit(token, payload, expiry)
	assert.NilError(t, err)

	got, exists, err := store.Find(token)
	assert.NilError(t, err)
	assert.Equal(t, exists, true)
	if !bytes.Equal(got, payload) {
		t.Errorf("got %q; want %q", got, payload)
	}

	updated := []byte("updated payload")
	err = store.Commit(token, updated, expiry)
	assert.NilError(t, err)

	got, exists, err = store.Find(token)
	assert.NilError(t, err)
	assert.Equal(t, exists, true)
	if !bytes.Equal(got, updated) {
		t.Errorf("got %q; want %q", got, updated)
	}

	err = store.Delete(token)
	assert.NilError(t, err)

	_, exists, err = store.Find(token)
	assert.NilError(t, err)
	assert.Equal(t, exists, false)
}

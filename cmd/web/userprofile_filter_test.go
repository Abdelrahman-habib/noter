package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Abdelrahman-habib/noter/internal/assert"
	"github.com/Abdelrahman-habib/noter/internal/models"
	"github.com/Abdelrahman-habib/noter/internal/models/mocks"
)

type recordingNoteModel struct {
	*mocks.NoteModel
	page      int
	limit     int
	public    *bool
	createdBy *int
}

func (m *recordingNoteModel) GetByPage(page int, limit int, public *bool, createdBy *int) ([]models.NoteWithUsername, models.PaginationMetaData, error) {
	m.page, m.limit, m.public, m.createdBy = page, limit, public, createdBy
	return m.NoteModel.GetByPage(page, limit, public, createdBy)
}

func newRecordingTestApplication(t *testing.T) (*application, *recordingNoteModel) {
	t.Helper()

	app := newTestApplication(t)
	notes := &recordingNoteModel{NoteModel: &mocks.NoteModel{}}
	app.notes = notes
	return app, notes
}

func loginAsAlice(t *testing.T, ts *testServer) {
	t.Helper()

	_, _, body := ts.get(t, "/user/login")
	form := url.Values{}
	form.Add("email", "alice@example.com")
	form.Add("password", "pa$$word")
	form.Add("csrf_token", extractCSRFToken(t, body))
	ts.postForm(t, "/user/login", form)
}

func assertShowFilter(t *testing.T, got, want *bool) {
	t.Helper()

	if (got == nil) != (want == nil) || (got != nil && *got != *want) {
		t.Errorf("show filter: got %s; want %s", showFilterString(got), showFilterString(want))
	}
}

func showFilterString(p *bool) string {
	switch {
	case p == nil:
		return "nil"
	case *p:
		return "&true"
	default:
		return "&false"
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestUserProfileShowFilter(t *testing.T) {
	tests := []struct {
		name          string
		asOwner       bool
		urlPath       string
		wantPublic    *bool
		wantCreatedBy int // the ID of the profile being viewed
	}{
		{
			name:          "Owner sees public notes",
			asOwner:       true,
			urlPath:       "/user/1?show=public",
			wantPublic:    boolPtr(true),
			wantCreatedBy: 1,
		},
		{
			name:          "Owner sees private notes",
			asOwner:       true,
			urlPath:       "/user/1?show=private",
			wantPublic:    boolPtr(false),
			wantCreatedBy: 1,
		},
		{
			name:          "Owner sees all notes",
			asOwner:       true,
			urlPath:       "/user/1?show=all",
			wantPublic:    nil,
			wantCreatedBy: 1,
		},
		{
			name:          "Owner without a show param sees all notes",
			asOwner:       true,
			urlPath:       "/user/1",
			wantPublic:    nil,
			wantCreatedBy: 1,
		},
		{
			name:          "Visitors see public notes even when asking for private ones",
			asOwner:       false,
			urlPath:       "/user/2?show=private",
			wantPublic:    boolPtr(true),
			wantCreatedBy: 2,
		},
		{
			name:          "Unauthenticated visitors see public notes",
			asOwner:       false,
			urlPath:       "/user/2",
			wantPublic:    boolPtr(true),
			wantCreatedBy: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, notes := newRecordingTestApplication(t)
			ts := newTestServer(t, app.routes())
			defer ts.Close()

			if tt.asOwner {
				loginAsAlice(t, ts)
			}

			code, _, _ := ts.get(t, tt.urlPath)
			assert.Equal(t, code, http.StatusOK)
			assertShowFilter(t, notes.public, tt.wantPublic)
			assert.Equal(t, notes.limit, 10)

			if notes.createdBy == nil {
				t.Errorf("createdBy: got nil; want %d", tt.wantCreatedBy)
			} else {
				assert.Equal(t, *notes.createdBy, tt.wantCreatedBy)
			}
		})
	}
}

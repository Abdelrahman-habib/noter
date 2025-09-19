package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Abdelrahman-habib/noter/internal/models"
	"github.com/Abdelrahman-habib/noter/ui"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear     int
	Note            models.NoteWithUsername
	Notes           []models.NoteWithUsername
	IsUserNote      bool
	NotesFilters    *models.NotesFilters
	User            models.User
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	CurrentPage     int
	HasNext         bool
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

// Helper function to safely check if a boolean pointer is not nil and true
func boolPtrIsTrue(b *bool) bool {
	return b != nil && *b
}

// Helper function to safely check if a boolean pointer is not nil and false
func boolPtrIsFalse(b *bool) bool {
	return b != nil && !*b
}

// Helper function to check if a boolean pointer is nil
func boolPtrIsNil(b *bool) bool {
	return b == nil
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"truncate":       truncate,
	"add":            add,
	"sub":            sub,
	"boolPtrIsTrue":  boolPtrIsTrue,
	"boolPtrIsFalse": boolPtrIsFalse,
	"boolPtrIsNil":   boolPtrIsNil,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := fs.Glob(ui.Files, "html/pages/*tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

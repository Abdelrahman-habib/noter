package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Abdelrahman-habib/noter/internal/models"
	"github.com/Abdelrahman-habib/noter/internal/validator"
	"github.com/google/uuid"
)

type noteUpsertForm struct {
	ID                  string `form:"id"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	Visibility          string `form:"visibility"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userChangePasswordForm struct {
	CurrentPassword     string `form:"currentPassword"`
	NewPassword         string `form:"newPassword"`
	ConfirmNewPassword  string `form:"confirmNewPassword"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	notes, err := app.notes.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Notes = notes

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) listNotes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	showPublic := true
	notes, metaData, err := app.notes.GetByPage(pageInt, 10, &showPublic, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.logger.Debug("meta data", slog.Any("metaData", metaData))

	data := app.newTemplateData(r)
	data.Notes = notes
	data.CurrentPage = pageInt
	data.HasNext = metaData.HasNext
	// Don't set NotesFilters for listNotes - this is for public notes only
	app.render(w, r, http.StatusOK, "list.tmpl", data)
}

func (app *application) myNotes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var showPublic *bool
	show, trueBool, falseBool := r.URL.Query().Get("show"), true, false

	if show != "" && show == "public" {
		showPublic = &trueBool
	} else if show == "private" {
		showPublic = &falseBool
	} // else leave it nil

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	notes, metaData, err := app.notes.GetByPage(pageInt, 10, showPublic, &userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.logger.Debug("meta data", slog.Any("metaData", metaData))

	data := app.newTemplateData(r)
	data.Notes = notes
	data.CurrentPage = pageInt
	data.HasNext = metaData.HasNext
	data.NotesFilters = &models.NotesFilters{
		ShowPublic: showPublic, // Pass the original pointer (can be nil, true, or false)
	}

	app.render(w, r, http.StatusOK, "list.tmpl", data)
}

func (app *application) noteView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := uuid.Validate(id)

	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	note, err := app.notes.Get(id, &userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Note = note
	data.IsUserNote = note.CreatedBy == userID

	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

func (app *application) noteCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = noteUpsertForm{
		ID:         "",
		Expires:    365,
		Visibility: "private",
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) noteCreatePost(w http.ResponseWriter, r *http.Request) {
	var form noteUpsertForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	form.CheckField(validator.PermittedValue(form.Visibility, "public", "private"), "visibility", "This field must equal public or private")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	noteID := form.ID
	isEditForm := form.ID != ""
	if isEditForm {
		err = uuid.Validate(noteID)
		if err != nil {
			app.clientError(w, http.StatusNotFound)
			return
		}
	}

	var id string
	if isEditForm {
		id, err = app.notes.Update(form.ID, form.Title, form.Content, form.Expires, form.Visibility == "public", userID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, http.StatusNotFound)
				return
			}
			app.serverError(w, r, err)
			return
		}
	} else {
		id, err = app.notes.Insert(form.Title, form.Content, form.Expires, form.Visibility == "public", userID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	flashMessage := "Note successfully created!"
	if isEditForm {
		flashMessage = "Note successfully updated!"
	}
	app.sessionManager.Put(r.Context(), "flash", flashMessage)

	// Redirect the user to the relevant page for the note.
	http.Redirect(w, r, fmt.Sprintf("/note/view/%v", id), http.StatusSeeOther)
}

func (app *application) noteEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)
	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	note, err := app.notes.Get(id, &userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, r, err)
		return
	}
	// expires is days between note.Created and note.Expires
	expires := int(note.Expires.Sub(note.Created).Hours() / 24)
	visibility := "private"
	if note.Public {
		visibility = "public"
	}
	data := app.newTemplateData(r)
	data.Form = noteUpsertForm{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Expires:    expires,
		Visibility: visibility,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) noteDeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)

	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.notes.Delete(id, &userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Note successfully deleted!")
	http.Redirect(w, r, "/my-notes", http.StatusSeeOther)
}

// Display a form for signing up a new user
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

// Create a new user
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "This field cannot be more than 255 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	_, err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Display a login form
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

// Log in a user
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("your email address or password is wrong")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/note/create", http.StatusSeeOther)
}

// log out a user
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfuly!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Display a login form
func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "change-password.tmpl", data)
}

// Log in a user
func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form userChangePasswordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.ConfirmNewPassword), "ConfirmNewPassword", "This field cannot be blank")
	form.CheckField(validator.EqualValue(form.ConfirmNewPassword, form.NewPassword), "confirmNewPassword", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "change-password.tmpl", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.ChangePassword(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "password is wrong")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "change-password.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about.tmpl", data)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.GetByID(userID)
	if err != nil {
		app.logger.Debug(err.Error())
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account.tmpl", data)

}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	user, err := app.users.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, r, err)
		return
	}

	// Check if the current user is viewing their own profile
	currentUserID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	isOwnProfile := currentUserID == id

	var showPublic *bool
	// If it's not their own profile, only show public notes
	if !isOwnProfile {
		trueBool := true
		showPublic = &trueBool
	}

	// Fetch notes for the user
	// We'll use a default page of 1, if no page is specified in the query string
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	notes, metaData, err := app.notes.GetByPage(pageInt, 10, showPublic, &id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	data.Notes = notes
	data.CurrentPage = pageInt
	data.HasNext = metaData.HasNext

	// Only pass NotesFilters if it's their own profile
	if isOwnProfile {
		// Parse 'show' query param to set active filter state
		show := r.URL.Query().Get("show")
		var filterShowPublic *bool
		switch show {
		case "public":
			trueBool := true
			filterShowPublic = &trueBool
		case "private":
			falseBool := false
			filterShowPublic = &falseBool
		}

		data.NotesFilters = &models.NotesFilters{
			ShowPublic: filterShowPublic,
		}
	}

	app.render(w, r, http.StatusOK, "profile.tmpl", data)
}

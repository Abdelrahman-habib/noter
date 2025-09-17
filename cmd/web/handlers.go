package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Abdelrahman-habib/snippetbox/internal/models"
	"github.com/Abdelrahman-habib/snippetbox/internal/validator"
	"github.com/google/uuid"
)

type snippetUpsertForm struct {
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
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) listSnippets(w http.ResponseWriter, r *http.Request) {
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
	snippets, metaData, err := app.snippets.GetByPage(pageInt, 10, &showPublic, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.logger.Debug("meta data", slog.Any("metaData", metaData))

	data := app.newTemplateData(r)
	data.Snippets = snippets
	data.CurrentPage = pageInt
	data.HasNext = metaData.HasNext
	// Don't set SnippetsFilters for listSnippets - this is for public snippets only
	app.render(w, r, http.StatusOK, "list.tmpl", data)
}

func (app *application) mySnippets(w http.ResponseWriter, r *http.Request) {
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

	snippets, metaData, err := app.snippets.GetByPage(pageInt, 10, showPublic, &userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.logger.Debug("meta data", slog.Any("metaData", metaData))

	data := app.newTemplateData(r)
	data.Snippets = snippets
	data.CurrentPage = pageInt
	data.HasNext = metaData.HasNext
	data.SnippetsFilters = &models.SnippetsFilters{
		ShowPublic: showPublic, // Pass the original pointer (can be nil, true, or false)
	}

	app.render(w, r, http.StatusOK, "list.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := uuid.Validate(id)

	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	snippet, err := app.snippets.Get(id, &userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.IsUserSnippet = snippet.CreatedBy == userID

	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetUpsertForm{
		ID:         "",
		Expires:    365,
		Visibility: "private",
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetUpsertForm

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
	snippetID := form.ID
	isEditForm := form.ID != ""
	if isEditForm {
		err = uuid.Validate(snippetID)
		if err != nil {
			app.clientError(w, http.StatusNotFound)
			return
		}
	}

	var id string
	if isEditForm {
		id, err = app.snippets.Update(form.ID, form.Title, form.Content, form.Expires, form.Visibility == "public", userID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, http.StatusNotFound)
				return
			}
			app.serverError(w, r, err)
			return
		}
	} else {
		id, err = app.snippets.Insert(form.Title, form.Content, form.Expires, form.Visibility == "public", userID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	flashMessage := "Snippet successfully created!"
	if isEditForm {
		flashMessage = "Snippet successfully updated!"
	}
	app.sessionManager.Put(r.Context(), "flash", flashMessage)

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)
}

func (app *application) snippetEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)
	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	snippet, err := app.snippets.Get(id, &userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, r, err)
		return
	}
	// expires is days between snippet.Created and snippet.Expires
	expires := int(snippet.Expires.Sub(snippet.Created).Hours() / 24)
	visibility := "private"
	if snippet.Public {
		visibility = "public"
	}
	data := app.newTemplateData(r)
	data.Form = snippetUpsertForm{
		ID:         snippet.ID,
		Title:      snippet.Title,
		Content:    snippet.Content,
		Expires:    expires,
		Visibility: visibility,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetDeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)

	if err != nil || id == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	app.logger.Debug("deleting snippet", slog.String("id", id))

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.snippets.Delete(id, &userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully deleted!")
	http.Redirect(w, r, "/my-snippets", http.StatusSeeOther)
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
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
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

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// ResetPassword render the password reset page.
// It takes HTTP response writer and request as parameters.
// It renders password reset page.
func (m *Repository) ResetPassword(w http.ResponseWriter, r *http.Request) {

	// create a data map that holds empty User model
	var emptyUser models.User
	data := make(map[string]interface{})
	data["reset_user"] = emptyUser

	// Render the "reset-password.page.tmpl" with form and data.
	render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// userStore the token and user data in cache
var userStore = &models.UserTokenStore{
	Users:             make(map[string]models.User),
	PasswordResetRepo: make(map[string]models.PasswordResetToken),
}

// PostResetPassword handles the post method that takes in the email address and send the reset token to user email.
// It takes HTTP response writer and request as parameters.
// It parse the form, validates the email, checks if email exists, then send reset token to email if exists.
func (m *Repository) PostResetPassword(w http.ResponseWriter, r *http.Request) {

	// Parse the form data from the request.
	// If any error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form with post form values.
	form := forms.New(r.PostForm)

	// Add form filed validation/
	form.Required("email")

	// Create a variable that stores the USER model/
	reset_user := models.User{
		Email: r.Form.Get("email"),
	}

	// Create a data map that holds the reset_user User model/
	data := make(map[string]interface{})
	data["reset_user"] = reset_user

	// Check if email exists.
	// If error occurs, a server error is returned.
	exists, err := m.DB.EmailExists(reset_user.Email)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If exists then add error to the form.
	if !exists {
		form.Errors.Add("email", "This email does not exist")
	}

	// If form is not valid then render the "reset-password.page.tmpl" go template with form and data.
	if !form.Valid() {
		render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// store the user email in cache.
	_ = userStore.Users[reset_user.Email]

	// create a token.
	token, err := helpers.GenerateRandomToken(32)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// Create a PasswordResetToken model and store token, email and expiry date
	resetToken := models.PasswordResetToken{
		Token:     token,
		Email:     reset_user.Email,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	userStore.PasswordResetRepo[token] = resetToken

	// send the email to email address with reset token
	body := fmt.Sprintf(`
		<h1>Reset Password</h1>
			<strong>The token for password change = </strong> %s<br><hr>
			<button><a href="%s/user/reset">Reset</a></button><br><hr>
			<strong>Ignore it, if you did not apply for reset password</strong>
	`, token, r.Header["Origin"][0])
	msg := models.MailData{
		To:      reset_user.Email,
		From:    m.App.AdminEmail,
		Subject: "Change Password",
		Content: body,
	}
	m.App.MailChan <- msg

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Reset token is sent to email")

	// Redirect user to password reset page.
	http.Redirect(w, r, "/user/reset", http.StatusSeeOther)
}

// ResetPasswordChange renders password change form with new and confirm password.
// It takes HTTP response writer and request as parameters.
func (m *Repository) ResetPasswordChange(w http.ResponseWriter, r *http.Request) {

	// Create a empty ResetPassword model
	var emptyPass models.ResetPassword

	// Create data map that holds emptyPass
	data := make(map[string]interface{})
	data["reset_password"] = emptyPass

	// render the "reset-password-change.page.tmpl" go template
	render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostResetPasswordChange handles post method and changes the password.
// It takes HTTP response writer and request as parameters.
// It parses the form data from the request, validates the form, checks the validity of the reset token, encrypts the new password,
// updates the user's password in the database, sends a notification email, and redirects the user to the login page.
func (m *Repository) PostResetPasswordChange(w http.ResponseWriter, r *http.Request) {

	// Parse the form data from the request.
	// If error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form wth parsed data
	form := forms.New(r.PostForm)
	passReset := models.ResetPassword{
		Token:              r.Form.Get("reset_token"),
		NewPassword:        r.Form.Get("new_password"),
		NewPasswordConfirm: r.Form.Get("confirm_new_password"),
	}

	// Form validation to form fileds
	form.Required("reset_token", "new_password", "confirm_new_password")

	// Create a data map to store teh reset_password
	data := make(map[string]interface{})
	data["reset_password"] = passReset

	// check if new password and confirm password matches
	if passReset.NewPassword != passReset.NewPasswordConfirm {
		form.Errors.Add("new_password", "password mismatch")
		form.Errors.Add("confirm_new_password", "password mismatch")
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// add password validations
	form.MinLength("new_password", 8)
	form.HasUpperCase("new_password")
	form.HasLowerCase("new_password")
	form.HasSpecialCharacter("new_password")
	form.HasNumber("new_password")

	// get the reset token
	resetToken := userStore.PasswordResetRepo[passReset.Token]
	if time.Now().After(resetToken.ExpiresAt) {
		form.Errors.Add("reset_token", "Token is invalid or expired")
	}

	// if form is not valid render the "reset-password-change.page.tmpl" with form and data
	if !form.Valid() {
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Encrypt the new password.
	hashed_password, err := helpers.EncryptPassword(passReset.NewPassword)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// Store the hashed password and email in user model.
	user := userStore.Users[resetToken.Email]
	user.Password = hashed_password
	userStore.Users[resetToken.Email] = user

	// Call ChangePassword interface to change the password.
	// If any error occurs, a server error is returned
	if err := m.DB.ChangePassword(hashed_password, resetToken.Email); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Deletes the token stored in cache.
	delete(userStore.PasswordResetRepo, resetToken.Token)
	// To do add notification for successfull password change
	body := fmt.Sprintf(`
		<h1>Password Reset Successfull<h1>
			<p>You can login from here<p><br>
			<button><a href="%s/user/login">Login</a></button>
	`, r.Header["Origin"][0])
	msg := models.MailData{
		To:      resetToken.Email,
		From:    m.App.AdminEmail,
		Subject: "Password Reset Successfull",
		Content: body,
	}

	// send the msg to email channel
	m.App.MailChan <- msg

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Password Reset Successfull")

	// Redirect user to login page if password reset is successfull
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

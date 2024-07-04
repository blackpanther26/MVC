package controllers

import (
	"net/http"
	"time"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/blackpanther26/mvc/pkg/views"
)

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	views.RenderTemplate(w, "signup", nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	views.RenderTemplate(w, "login", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			views.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Failed to parse form"})
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if !isPasswordValid(password) && len(password) > 0 {
			views.RenderTemplate(w, "signup", map[string]interface{}{
				"ErrorMessage": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one digit",
			})
			return
		}

		if password != confirmPassword {
			views.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Passwords do not match"})
			return
		}

		err = models.CreateUser(username, password)
		if err != nil {
			views.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Failed to create user"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		views.RenderTemplate(w, "signup", nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		views.RenderTemplate(w, "login", map[string]interface{}{"ErrorMessage": "Failed to parse form"})
		return
	}

	body.Username = r.FormValue("username")
	body.Password = r.FormValue("password")

	user, err := models.AuthenticateUser(body.Username, body.Password)
	if err != nil {
		views.RenderTemplate(w, "login", map[string]interface{}{"ErrorMessage": "Invalid username or password"})
		return
	}

	token, err := models.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Authorization", Value: token, Expires: time.Now().Add(24 * time.Hour)}
	http.SetCookie(w, cookie)

	if user.IsAdmin {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/client/", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour * 24),
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func isPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpperCase := false
	hasLowerCase := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpperCase = true
		case 'a' <= char && char <= 'z':
			hasLowerCase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	return hasUpperCase && hasLowerCase && hasDigit
}
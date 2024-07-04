package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/blackpanther26/mvc/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "signup", nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			utils.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Failed to parse form"})
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if !isPasswordValid(password) && len(password) > 0 {
			utils.RenderTemplate(w, "signup", map[string]interface{}{
				"ErrorMessage": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one digit",
			})
			return
		}

		if password != confirmPassword {
			utils.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Passwords do not match"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			utils.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Failed to hash password"})
			return
		}

		user := models.User{
			Username:     username,
			PasswordHash: string(hash),
		}

		result := config.DB.Create(&user)
		if result.Error != nil {
			utils.RenderTemplate(w, "signup", map[string]interface{}{"ErrorMessage": "Failed to create user"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		utils.RenderTemplate(w, "signup", nil)
	}
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		utils.RenderTemplate(w, "login", map[string]interface{}{"ErrorMessage": "Failed to parse form"})
		return
	}

	body.Username = r.FormValue("username")
	body.Password = r.FormValue("password")

	var user models.User
	config.DB.First(&user, "username =?", body.Username)

	if user.ID == 0 {
		utils.RenderTemplate(w, "login", map[string]interface{}{"ErrorMessage": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		utils.RenderTemplate(w, "login", map[string]interface{}{"ErrorMessage": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("SECRET")
	if secret == "" {
		http.Error(w, "SECRET environment variable not set", http.StatusInternalServerError)
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Authorization", Value: tokenString, Expires: time.Now().Add(24 * time.Hour)}
	http.SetCookie(w, cookie)

	if user.IsAdmin {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/client/", http.StatusSeeOther)
	}
}

// func Validate(w http.ResponseWriter, r *http.Request) {
// 	userCtx := r.Context()
// 	user := userCtx.Value("user")

// 	userJson, err := json.Marshal(user)
// 	if err!= nil {
// 		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(userJson)
// }

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
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/blackpanther26/mvc/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi from signup handler")
	utils.RenderTemplate(w, "signup", nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read body"}`, http.StatusBadRequest)
		return
	}

	if body.Username == "" || body.Password == "" {
		http.Error(w, `{"error":"Username and Password are required"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     body.Username,
		PasswordHash: string(hash),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, `{"error":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Signup successful", "user_id":` + strconv.Itoa(int(user.ID)) + `}`))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var user models.User
	config.DB.First(&user, "username =?", body.Username)

	if user.ID == 0 {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Login successful"}`))
}

func Validate(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context()
	user := userCtx.Value("user")

	userJson, err := json.Marshal(user)
	if err!= nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

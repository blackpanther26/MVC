package models

import (
	"errors"
	"os"
	"time"

	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var userCount int64
	if err := config.DB.Model(&types.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	user := types.User{
		Username:     username,
		PasswordHash: string(hash),
		IsAdmin:      userCount == 0,
	}

	result := config.DB.Create(&user)
	return result.Error
}

func AuthenticateUser(username, password string) (*types.User, error) {
	var user types.User
	result := config.DB.First(&user, "username =?", username)
	if result.Error != nil {
		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GenerateToken(userID uint) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", errors.New("SECRET environment variable not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

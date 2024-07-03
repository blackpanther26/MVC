package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("Authorization")
		if err != nil || tokenCookie == nil || tokenCookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		secret := []byte(os.Getenv("SECRET"))

		token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			currentTime := float64(time.Now().Unix())

			if expClaim, exists := claims["exp"]; exists {
				if currentTime > expClaim.(float64) {
					log.Println("Token has expired")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			var user models.User
			config.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				log.Println("User not found")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Println("No claims found in token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func WithRequireAuth(next http.Handler) http.Handler {
    return RequireAuth(next)
}

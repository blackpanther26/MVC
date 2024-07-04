package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/types"
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
			fmt.Printf("Error parsing token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Println("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			currentTime := float64(time.Now().Unix())
			if expClaim, exists := claims["exp"]; exists {
				if currentTime > expClaim.(float64) {
					fmt.Println("Token has expired")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			var user types.User
			config.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("No claims found in token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
package middleware

import (
	"net/http"
	"github.com/blackpanther26/mvc/pkg/types"
)

func IsNotAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        userCtx := r.Context().Value("user")
        if userCtx == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        user, ok := userCtx.(types.User)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if user.IsAdmin {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

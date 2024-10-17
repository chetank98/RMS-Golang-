package Middleware

import (
	"RMS/Models"
	"net/http"
)

func ShouldHaveRole(role Models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := UserContext(r).Role
			if userRole != role {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

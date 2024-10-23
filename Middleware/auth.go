package Middleware

import (
	"RMS/Database/DbHelper"
	"RMS/Models"
	"RMS/Utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

type ContextKey string

const (
	userContext ContextKey = "userContext"
)

//func Authenticate(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			Utils.RespondError(w, http.StatusUnauthorized, nil, "authorization header missing")
//			return
//		}
//
//		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//		if tokenString == authHeader {
//			Utils.RespondError(w, http.StatusUnauthorized, nil, "bearer token missing")
//			return
//		}
//		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, errors.New("invalid signing method")
//			}
//			return []byte(os.Getenv("SECRET_KEY")), nil
//		})
//
//		if parseErr != nil || !token.Valid {
//			Utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
//			return
//		}
//
//		claimValues, ok := token.Claims.(jwt.MapClaims)
//		if !ok || !token.Valid {
//			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
//			return
//		}
//
//		sessionID := claimValues["sessionId"].(string)
//		archivedAt, err := DbHelper.GetArchivedAt(sessionID)
//		if err != nil {
//			Utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
//			return
//		}
//
//		if archivedAt != nil {
//			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token")
//			return
//		}
//
//		user := &Models.UserCtx{
//			UserID:    claimValues["userId"].(string),
//			SessionID: sessionID,
//			Role:      Models.Role(claimValues["role"].(string)),
//		}
//
//		ctx := context.WithValue(r.Context(), userContext, user)
//		r = r.WithContext(ctx)
//
//		next.ServeHTTP(w, r)
//
//	})
//}

func SetCookies(w http.ResponseWriter, tokenString string) {
	cookie := http.Cookie{
		Name:     "access-token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		at, err := r.Cookie("access-token")
		if err != nil {
			if err == http.ErrNoCookie {
				Utils.RespondError(w, http.StatusUnauthorized, nil, "authorization cookie missing")
			} else {
				Utils.RespondError(w, http.StatusBadRequest, err, "error retrieving cookie")
			}
			return
		}

		tokenString := at.Value

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method") // Invalid signing method error
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if parseErr != nil || !token.Valid {
			Utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid claims")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token ")
			return
		}

		sessionID := claimValues["sessionId"].(string)

		archivedAt, err := DbHelper.GetArchivedAt(sessionID)
		if err != nil {
			Utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
			return
		}

		if archivedAt != nil {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token")
			return
		}

		user := &Models.UserCtx{
			UserID:    claimValues["userId"].(string),
			SessionID: sessionID,
			Role:      Models.Role(claimValues["role"].(string)),
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserContext(r *http.Request) *Models.UserCtx {
	if user, ok := r.Context().Value(userContext).(*Models.UserCtx); ok {
		return user
	}
	return nil
}

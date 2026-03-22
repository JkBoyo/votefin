package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/internal/jellyfin"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) logoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		return
	}

	cookie.MaxAge = -10
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
		return
	}

	userName := r.FormValue("Username")

	passWord := r.FormValue("Password")

	authResp, err := jellyfin.AuthenticateUser(userName, passWord, r.Context())
	if err == jellyfin.JellyfinAuthError {
		respondWithHTML(w, http.StatusAccepted, templates.Notification("Invalid username or password"))
		slog.Warn("invalid username or password entered")
		return
	} else if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, "Error setting authentication")
		return
	}
	_, err = cfg.db.GetUserByJellyID(r.Context(), authResp.User.Id)
	if err == sql.ErrNoRows {
		currTime := time.Now().Unix()
		var isAdmin int64
		if authResp.User.Policy.IsAdministrator {
			isAdmin = 1

		} else {
			isAdmin = 0
		}
		newUser := database.AddUserParams{
			CreatedAt:      currTime,
			UpdatedAt:      currTime,
			JellyfinUserID: authResp.User.Id,
			Username:       authResp.User.Name,
			IsAdmin:        isAdmin,
		}

		_, err = cfg.db.AddUser(r.Context(), newUser)
		if err != nil {
			fmt.Println()
			fmt.Println("error adding user" + err.Error())
		}
	}
	// INFO: Sets the cookie to require tls or not only allowing local login when false. Defaults to true
	isCookieSecureStr := os.Getenv("SECURE_ONLY_COOKIE")
	if isCookieSecureStr == "" {
		isCookieSecureStr = "true"
	}
	isCookieSecure, err := strconv.ParseBool(isCookieSecureStr)
	if err != nil {
		slog.Error("Error converting secure cookie flag to a bool", "error", err)
		return
	}

	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    authResp.AccessToken,
		Secure:   isCookieSecure,
		HttpOnly: true,
		SameSite: 2,
		Path:     "/",
	}

	http.SetCookie(w, authCookie)

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

type ContextKey string

const (
	userContextKey ContextKey = "user"
)

func (cfg *apiConfig) AuthorizeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Cookies())
		cookie, err := r.Cookie("Token")
		if err != nil {
			next.ServeHTTP(w, r)
			fmt.Println("Error getting cookie: " + err.Error())
			return
		}
		token := cookie.Value
		if token == "" {
			fmt.Println("No token in the cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		jfUser, err := jellyfin.ValidateToken(token)
		if err != nil {
			fmt.Println("Error validating token")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := cfg.db.GetUserByJellyID(r.Context(), jfUser.Id)
		if err == sql.ErrNoRows {
			currTime := time.Now().Unix()
			var isAdmin int64
			if jfUser.Policy.IsAdministrator {
				isAdmin = 1
			} else {
				isAdmin = 0
			}
			newUser := database.AddUserParams{
				CreatedAt:      currTime,
				UpdatedAt:      currTime,
				JellyfinUserID: jfUser.Id,
				Username:       jfUser.Name,
				IsAdmin:        isAdmin,
			}
			fmt.Println(newUser)
			user, err = cfg.db.AddUser(r.Context(), newUser)
			if err != nil {
				fmt.Println("error adding user: ", err.Error())
			}
		}

		con := context.WithValue(r.Context(), userContextKey, &user)

		next.ServeHTTP(w, r.WithContext(con))
	},
	)
}

func CheckIsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hitting admin endpoint")
		fmt.Println(r.Context())
		user, ok := r.Context().Value(userContextKey).(*database.User)
		fmt.Println(user)
		if !ok {
			respondWithHtmlErr(w, http.StatusInternalServerError, "User not found")
			slog.Error("error user not found in request context")
			return
		}

		if user.IsAdmin == 0 {
			respondWithHtmlErr(w, http.StatusForbidden, "User unauthorized")
			slog.Warn("Non admin user hitting admin endpoint", "username", user.Username, "userid", user.JellyfinUserID)
			return
		}

		next.ServeHTTP(w, r)
	})
}

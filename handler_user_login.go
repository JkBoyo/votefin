package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/internal/jellyfin"
	"www.github.com/jkboyo/votefin/templates"
)

type authorizedHandler func(w http.ResponseWriter, r *http.Request, user *database.User)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}

	userName := r.FormValue("Username")

	passWord := r.FormValue("Password")

	authResp, err := jellyfin.AuthenticateUser(userName, passWord, r.Context())
	if err == jellyfin.JellyfinAuthError {
		respondWithHtmlErr(w, http.StatusUnauthorized, "Invalid username or password")
	} else if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, "Error setting authentication")
	}
	_, err = cfg.db.GetUserByJellyID(r.Context(), authResp.User.Id)
	if err == sql.ErrNoRows {
		currTime := time.Now().Local().String()
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

	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    authResp.AccessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	}

	http.SetCookie(w, authCookie)

	http.Redirect(w, r, "/dashboard", http.StatusAccepted)
}

func (cfg *apiConfig) AuthorizeHandler(handler authorizedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		if err != nil {
			respondWithHTML(w, http.StatusAccepted, templates.Notification(err.Error()))
			handler(w, r, nil)
		}
		token := cookie.Value
		if token == "" {
			fmt.Println("No token in the cookie")
		}

		jfUser, err := jellyfin.ValidateToken(token)
		if err != nil {
			fmt.Println("Error validating token")
			respondWithHtmlErr(w, http.StatusUnauthorized, ("Error authenticating user: " + err.Error()))
		}

		user, err := cfg.db.GetUserByJellyID(r.Context(), jfUser.Id)
		if err == sql.ErrNoRows {
			currTime := time.Now().Local().String()
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

		handler(w, r, &user)
	}
}

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

type authorizedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}

	userName := r.FormValue("Username")

	passWord := r.FormValue("Password")

	fmt.Println("User Name: " + userName)
	fmt.Println("Password: " + passWord)

	authResp, err := jellyfin.AuthenticateUser(userName, passWord, r.Context())
	if err == jellyfin.JellyfinAuthError {
		respondWithHTML(w, http.StatusUnauthorized, (templates.LoginError("Invalid username or password")))
	} else if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, "Error setting authentication")
	}
	user, err := cfg.db.GetUserByJellyID(r.Context(), authResp.User.Id)
	if err == sql.ErrNoRows {
		currTime := time.Now().Local().String()
		var isAdmin int64
		if authResp.User.IsAdmin {
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
		fmt.Println(newUser)
		user, err = cfg.db.AddUser(r.Context(), newUser)
	}

	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    authResp.AccessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	}

	http.SetCookie(w, authCookie)

	page, err := renderPage(cfg, r, user)
	if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
	}

	respondWithHTML(w, http.StatusAccepted, page)
}

func (cfg *apiConfig) AuthorizeHandler(handler authorizedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			respondWithHTML(w, http.StatusAccepted, templates.LoginError(err.Error()))
			return
		}
		token := cookie.Value

		jfUser, err := jellyfin.ValidateToken(token)

		user, err := cfg.db.GetUserByJellyID(r.Context(), jfUser.Id)
		if err == sql.ErrNoRows {
			currTime := time.Now().Local().String()
			var isAdmin int64
			if jfUser.IsAdmin {
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
		}

		handler(w, r, user)
	}
}

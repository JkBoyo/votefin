package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/jellyfin"
	"www.github.com/jkboyo/votefin/templates"
)

type authorizedHandler func(w http.ResponseWriter, r *http.Request, user jellyfin.JellyfinUser)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}

	userName := r.FormValue("userName")

	passWord := r.FormValue("passWord")

	fmt.Println("User Name: " + userName)
	fmt.Println("Password: " + passWord)

	authResp, err := jellyfin.AuthenticateUser(userName, passWord, r.Context())

	if err == jellyfin.JellyfinAuthError {

		respondWithHTML(w, http.StatusUnauthorized, templates.BasePage(templates.Login()))
	} else if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, "Error setting authentication")
	}

	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    authResp.AccessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	}

	http.SetCookie(w, authCookie)

	respondWithHTML(w, http.StatusAccepted, mainPage)
}

func (api *apiConfig) AuthorizeHandler(handler authorizedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			respondWithHTML(w, http.StatusNetworkAuthenticationRequired, templates.BasePage(templates.Login()))
			return
		}
		token := cookie.Value

		user, err := jellyfin.ValidateToken(token)

		handler(w, r, user)
	}
}

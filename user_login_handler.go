package main

import (
	"fmt"
	"net/http"

	"www.github.com/jkboyo/votefin/internal/jellyfin"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}

	userName := r.FormValue("userName")

	passWord := r.FormValue("passWord")

	fmt.Println("User Name: " + userName)
	fmt.Println("Password: " + passWord)

	token, err := jellyfin.AuthenticateUser(userName, passWord, r.Context())

	if err == jellyfin.JellyfinAuthError {

		respondWithHTML(w, http.StatusUnauthorized, templates.BasePage(templates.Login()))
	} else if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, "Error setting authentication")
	}

	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	}

	http.SetCookie(w, authCookie)

	mainPage := templates.BasePage(templates.VotePage()) //TODO: Fill in elements of the main page

	respondWithHTML(w, http.StatusAccepted, mainPage)
}

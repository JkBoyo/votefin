package main

import (
	"fmt"
	"net/http"

	"www.github.com/jkboyo/votefin/internal/jellyfin"
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

	if err != nil {
		fmt.Printf("Authentication failed: %s\n", err)
	}
	authCookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	}

	http.SetCookie(w, authCookie)
}

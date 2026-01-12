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

func respondWithHTML(w http.ResponseWriter, code int, comp templ.Component) error {
	w.Header().Set("Content-Type", "application-/x-www-form-urlencoded")
	w.WriteHeader(code)
	err := templates.BasePage(comp).Render(context.Background(), os.Stdout)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
	}
	err = templates.BasePage(comp).Render(context.Background(), w)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
	}
	return nil
}

func respondWithHtmlErr(w http.ResponseWriter, code int, errMsg string) error {
	errNotif := templates.Notification(errMsg)
	return respondWithHTML(w, code, errNotif)
}

func (cfg *apiConfig) renderPageHandler(w http.ResponseWriter, r *http.Request, jfUser jellyfin.JellyfinUser) {
	page, err := renderPage(cfg, r, jfUser)
	if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
	}

	respondWithHTML(w, http.StatusAccepted, page)
}

func renderPage(cfg *apiConfig, r *http.Request, u jellyfin.JellyfinUser) (templ.Component, error) {
	votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on: %v", err.Error())
	}

	user, err := cfg.db.GetUserByJellyID(r.Context(), u.Id)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user info: %v", err.Error())
	}

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on by the user: %v", err.Error())
	}

	allMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching all movies for voting: %v", err.Error())
	}

	return templates.VotePage(votedOnMovies, userVotedMovies, allMovies), nil
}

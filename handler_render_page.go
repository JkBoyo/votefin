package main

import (
	"context"
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
	err = comp.Render(context.Background(), w)
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
	votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
	if err != nil {
		//TODO: set error handling
	}

	user, err := cfg.db.GetUserByJellyID(r.Context(), jfUser.Id)

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), user.ID)
	if err != nil {
		//TODO: set error handling
	}

	allMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		//TODO: set error handling
	}

	page := templates.BasePage(templates.VotePage(votedOnMovies, userVotedMovies, allMovies))
}

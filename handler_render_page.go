package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

func respondWithHTML(w http.ResponseWriter, code int, comp templ.Component) error {
	w.Header().Set("Content-Type", "application-/x-www-form-urlencoded")
	w.WriteHeader(code)
	err := comp.Render(context.Background(), os.Stdout)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
		return err
	}
	err = comp.Render(context.Background(), w)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
		return err
	}
	return nil
}

func respondWithHtmlErr(w http.ResponseWriter, code int, errMsg string) error {
	errNotif := templates.Notification(errMsg)
	return respondWithHTML(w, code, errNotif)
}

func (cfg *apiConfig) renderPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		slog.Error("couldn't access user from context for dashboard")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	page, err := renderPage(cfg, r, user)
	if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	page.Render(r.Context(), w)
}

func renderPage(cfg *apiConfig, r *http.Request, user *database.User) (templ.Component, error) {
	votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on: %v", err.Error())
	}

	userVotesCount, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err != nil {
		fmt.Println("get votecount failed: " + err.Error())
		return nil, fmt.Errorf("Error fetching user vote count: %v", err.Error())
	}

	userVotesLeft := cfg.voteLimit - int(userVotesCount.Float64)

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on by the user: %v", err.Error())
	}

	allMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching all movies for voting: %v", err.Error())
	}

	return templates.BasePage(templates.VotePage(user, user.IsAdmin, votedOnMovies, userVotesLeft, userVotedMovies, allMovies)), nil
}

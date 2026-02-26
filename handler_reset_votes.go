package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"www.github.com/jkboyo/votefin/internal/database"
)

type movieStatus string

const (
	statusOnServer   movieStatus = "on server"
	statusProcessing movieStatus = "processing"
	statusNotRipped  movieStatus = "not ripped"
)

func (cfg *apiConfig) removeMovie(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Error parsing form data", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	movieId, err := strconv.ParseInt(r.URL.Query().Get("movieId"), 10, 64)
	if err != nil {
		slog.Error("Error converting query results to integer", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't convert movie ID to integer")
		return
	}

	err = cfg.db.DeleteMovie(r.Context(), movieId)
	if err != nil {
		slog.Error("Error removing movie from database", "error", err)
	}

	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		slog.Error("couldn't access user from context for mark finished endpoint")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
		return
	}

	updateVotesPage(cfg, r, w, user, currUserVotes)
}
func (cfg *apiConfig) markFinishedhandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Error parsing form data", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	movieId, err := strconv.ParseInt(r.URL.Query().Get("movieId"), 10, 64)
	if err != nil {
		slog.Error("Error converting query results to integer", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't convert movie ID to integer")
		return
	}

	err = resetVotes(r, cfg, movieId)
	if err != nil {
		slog.Error("Error reseting votes to zero", "error", err)
		return
	}

	setZeroArgs := database.SetMovieVotesZeroParams{
		UpdatedAt: time.Now().Unix(),
		MovieID:   movieId,
	}
	err = cfg.db.SetMovieVotesZero(r.Context(), setZeroArgs)
	if err != nil {
		slog.Error("Error setting movie votes to zero", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, "couldn't set votes to zero")
		return
	}

	setStatusArgs := database.UpdateStatusParams{
		UpdatedAt: time.Now().Unix(),
		Status:    string(statusOnServer),
		ID:        movieId,
	}
	err = cfg.db.UpdateStatus(r.Context(), setStatusArgs)
	if err != nil {
		slog.Error("Error setting status", "error", err)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		slog.Error("couldn't access user from context for mark finished endpoint")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
		return
	}

	updateVotesPage(cfg, r, w, user, currUserVotes)
}

func resetVotes(r *http.Request, cfg *apiConfig, movie int64) error {
	args := database.SetMovieVotesZeroParams{
		UpdatedAt: time.Now().Unix(),
		MovieID:   movie,
	}
	err := cfg.db.SetMovieVotesZero(r.Context(), args)
	if err != nil {
		return err
	}
	return nil
}

package main

import (
	"net/http"

	"www.github.com/jkboyo/votefin/internal/jellyfin"
)

func (cfg *apiConfig) renderPageHandler(w http.ResponseWriter, r *http.Request, user jellyfin.JellyfinUser) {
	allMovies, err := cfg.db.GetMovies(r.Context())

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context())
}

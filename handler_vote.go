package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) voteHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		respondWithHtmlErr(w, http.StatusUnauthorized, "Not authorized to vote")
		return
	}
	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
	}

	movieId := r.URL.Query().Get("movieId")
	if err != nil {
		slog.Error("Error converting query results to integer", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't convert movie ID to integer")
		return
	}

	if int(currUserVotes.Float64) == cfg.voteLimit {
		respondWithHtmlErr(w,
			http.StatusAccepted,
			"Too Many votes submitted. You have no votes left",
		)
		return
	}

	currUserVotes, err = updateVotes(r, user, cfg, currUserVotes, movieId, +1)

	updateVotesPage(cfg, r, w, user, currUserVotes)
}

func (cfg *apiConfig) voteRemovalHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		respondWithHtmlErr(w, http.StatusUnauthorized, "Not authorized to vote")
		return
	}

	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
	}

	movieId := r.URL.Query().Get("movieId")
	if err != nil {
		slog.Error("Error converting query results to integer", "error", err)
		respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't convert movie ID to integer")
		return
	}

	if currUserVotes.Float64 == 0 { // Should Never happen
		respondWithHtmlErr(w,
			http.StatusForbidden,
			"No votes currently had can't go negative",
		)
		slog.Error("Can't have negative votes")
		return
	}

	currUserVotes, err = updateVotes(r, user, cfg, currUserVotes, movieId, -1)

	updateVotesPage(cfg, r, w, user, currUserVotes)
}

func updateVotes(r *http.Request, user *database.User, cfg *apiConfig, currUserVotes sql.NullFloat64, movie string, change int64) (sql.NullFloat64, error) {
	movieId, err := strconv.ParseInt(movie, 10, 64)
	if err != nil {
		slog.Error("Error parsing movie id", "error", err)
		return sql.NullFloat64{}, err
	}
	checkMovie := database.CheckMovieVotesParams{
		UserID:  user.ID,
		MovieID: movieId,
	}
	vote, err := cfg.db.CheckMovieVotes(r.Context(), checkMovie)
	if err == sql.ErrNoRows { // Handle no vote being present
		voteParams := database.CreateVoteParams{
			CreatedAt: time.Now().Unix(),
			UserID:    user.ID,
			MovieID:   movieId,
			VoteCount: int64(change),
		}
		_, err = cfg.db.CreateVote(r.Context(), voteParams)
		if err != nil {
			slog.Error("Error creating vote from vote handler", "error", err)
			return sql.NullFloat64{}, err
		}

		currUserVotes.Float64 += float64(change)
		return currUserVotes, nil
	} else if err != nil {
		slog.Error("Error checking on vote based on user and movie id's", "error", err)
		return sql.NullFloat64{}, err
	} else {
		updateVote := database.UpdateVoteCountParams{
			UpdatedAt: time.Now().Unix(),
			VoteCount: vote.VoteCount + change,
			ID:        vote.ID,
		}

		err = cfg.db.UpdateVoteCount(r.Context(), updateVote)
		if err != nil {
			slog.Error("Error updating votecount to add more votes", "error", err)
			return sql.NullFloat64{}, err
		}
		currUserVotes.Float64 += float64(change)
		return currUserVotes, nil
	}
}

func updateVotesPage(cfg *apiConfig, r *http.Request, w http.ResponseWriter, user *database.User, currUserVotes sql.NullFloat64) {
	votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
	if err != nil {
		slog.Error("Error getting movies sorted by votes", "error", err)
		return
	}

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), user.ID)
	if err != nil {
		slog.Error("Error getting user voted movies sorted by votes", "error", err)
		return
	}

	allMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		slog.Error("Error getting user voted movies sorted by votes", "error", err)
		return
	}

	respondWithHTML(w, http.StatusAccepted, templ.Join(
		templates.VotesMovieList(user.IsAdmin, true, votedOnMovies),
		templates.UserVotesMovieList(true, userVotedMovies),
		templates.MovieList(user.IsAdmin, cfg.voteLimit-int(currUserVotes.Float64), true, allMovies),
		templates.Notification(""),
	),
	)
}

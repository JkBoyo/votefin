package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) voteHandler(w http.ResponseWriter, r *http.Request, user *database.User) {
	if user == nil {
		respondWithHtmlErr(w, http.StatusUnauthorized, "Not authorized to vote")
		return
	}
	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error Parsing votes")
		respondWithHtmlErr(w, http.StatusBadRequest, "Error parsing vote request"+err.Error())
	}

	fmt.Println(r.Form)

	votedMovies := r.PostForm["votedMovie"]

	if len(votedMovies)+int(currUserVotes.Float64) > cfg.voteLimit { //TODO: when implementing roles apply role based multiplier
		respondWithHtmlErr(w,
			http.StatusForbidden,
			fmt.Sprintf("Too Many votes submitted. You have %d votes left", cfg.voteLimit-int(currUserVotes.Float64)),
		)
		return
	}

	for _, movie := range votedMovies {
		currUserVotes, err = updateVotes(w, r, user, cfg, currUserVotes, movie, +1)
	}

	updateVotesPage(cfg, r, w, user, currUserVotes)
}

func (cfg *apiConfig) voteRemovalHandler(w http.ResponseWriter, r *http.Request, user *database.User) {
	if user == nil {
		respondWithHtmlErr(w, http.StatusUnauthorized, "Not authorized to vote")
		return
	}

	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		slog.Warn("No rows returned", "error", err)
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error Parsing votes")
		respondWithHtmlErr(w, http.StatusBadRequest, "Error parsing vote request"+err.Error())
	}

	fmt.Println(r.Form)

	unVotedMovies := r.PostForm["votedMovie"]

	if int(currUserVotes.Float64)-len(unVotedMovies) < 0 { // Should Never happen
		respondWithHtmlErr(w,
			http.StatusForbidden,
			fmt.Sprintf("Too Many votes submitted. You have %d votes left", cfg.voteLimit-int(currUserVotes.Float64)),
		)
		slog.Error("Can't have negative votes")
		return
	}

	for _, movie := range unVotedMovies {
		currUserVotes, err = updateVotes(w, r, user, cfg, currUserVotes, movie, -1)
	}

	updateVotesPage(cfg, r, w, user, currUserVotes)
}

func updateVotes(w http.ResponseWriter, r *http.Request, user *database.User, cfg *apiConfig, currUserVotes sql.NullFloat64, movie string, change int64) (sql.NullFloat64, error) {
	movieId, err := strconv.ParseInt(movie, 10, 64)
	if err != nil {
		fmt.Println("Couldn't convert", movie, "to an integer")
		respondWithHtmlErr(w, http.StatusBadRequest, "Couldn't convert movie id to integer")
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
			respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't insert vote into db"+err.Error())
			slog.Error("Error creating vote from vote handler", "error", err)
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

	respondWithHTML(w, http.StatusAccepted, templ.Join(
		templates.VotesMovieList(true, "moviesVotedOn", votedOnMovies),
		templates.UserVotesMovieList(true, "userMoviesVotedOn", cfg.voteLimit-int(currUserVotes.Float64), userVotedMovies),
	),
	)
}

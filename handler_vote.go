package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) voteHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), u.ID)
	if err == sql.ErrNoRows {

	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error Parsing votes")
		respondWithHtmlErr(w, http.StatusBadRequest, "Error parsing vote request")
	}

	votedMovies := r.PostForm["votedMovies"]

	voteLim := cfg.voteLimit

	if len(votedMovies)+int(currUserVotes) > voteLim { //TODO: when implementing roles apply role based multiplier
		respondWithHtmlErr(w,
			http.StatusForbidden,
			fmt.Sprintf("Too Many votes submitted. You have %d votes left", voteLim-int(currUserVotes)),
		)
	}

	for _, movie := range votedMovies {
		movieId, err := strconv.ParseInt(movie, 10, 64)
		if err != nil {
			fmt.Println("Couldn't convert", movie, "to an integer")
			respondWithHtmlErr(w, http.StatusBadRequest, "Couldn't convert movie id to integer")
		}
		voteParams := database.CreateVoteParams{
			CreatedAt: time.Now().Unix(),
			UserID:    u.ID,
			MovieID:   movieId,
		}
		_, err = cfg.db.CreateVote(r.Context(), voteParams)

	}

	respondWithHTML(w, http.StatusAccepted)
}

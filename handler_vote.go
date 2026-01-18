package main

import (
	"net/http"

	"www.github.com/jkboyo/votefin/internal/database"
)

func (cfg *apiConfig) voteHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), u.ID)
	if err != nil {
		//TODO: handle this error
	}

	err = r.ParseForm()
	if err != nil {
		//TODO: handle form parsing error
	}

	votedMovies := r.PostForm["votedMovies"]
	if len(votedMovies)+int(currUserVotes) > voteLimit {

	}
}

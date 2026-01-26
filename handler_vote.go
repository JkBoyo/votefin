package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) voteHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	currUserVotes, err := cfg.db.GetVotesCountPerUser(r.Context(), u.ID)
	if err == sql.ErrNoRows {
		currUserVotes = 0
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error Parsing votes")
		respondWithHtmlErr(w, http.StatusBadRequest, "Error parsing vote request"+err.Error())
	}

	fmt.Println(r.Form)

	votedMovies := r.PostForm["votedMovie"]

	voteLim := cfg.voteLimit

	if len(votedMovies)+int(currUserVotes) > voteLim { //TODO: when implementing roles apply role based multiplier
		respondWithHtmlErr(w,
			http.StatusForbidden,
			fmt.Sprintf("Too Many votes submitted. You have %d votes left", voteLim-int(currUserVotes)),
		)
	} else {
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
			if err != nil {
				respondWithHtmlErr(w, http.StatusInternalServerError, "Couldn't insert vote into db"+err.Error())
			}

		}

		votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
		if err != nil {
			respondWithHtmlErr(w, http.StatusInternalServerError, "Error fetching voted on movies: "+err.Error())
		}

		userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), u.ID)
		if err != nil {
			respondWithHtmlErr(w, http.StatusInternalServerError, "Error fetching users voted on movies: "+err.Error())
		}

		respondWithHTML(w, http.StatusAccepted, templ.Join(
			templates.VotesMovieList(true, "moviesVotedOn", votedOnMovies),
			templates.UserVotesMovieList(true, "userMoviesVotedOn", userVotedMovies),
		),
		)
	}

}

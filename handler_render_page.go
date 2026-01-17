package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"www.github.com/jkboyo/votefin/internal/database"
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

func (cfg *apiConfig) renderPageHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	page, err := renderPage(cfg, r, u)
	if err != nil {
		respondWithHtmlErr(w, http.StatusInternalServerError, err.Error())
	}

	respondWithHTML(w, http.StatusAccepted, page)
}

func renderPage(cfg *apiConfig, r *http.Request, u database.User) (templ.Component, error) {
	votedOnMovies, err := cfg.db.GetMoviesSortedByVotes(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on: %v", err.Error())
	}

	fmt.Println(u)

	user, err := cfg.db.GetUserByJellyID(r.Context(), u.Id)
	if err == sql.ErrNoRows {
		currTime := time.Now().Local().String()
		var isAdmin int64
		if u.IsAdmin {
			isAdmin = 1
		} else {
			isAdmin = 0
		}
		newUser := database.AddUserParams{
			CreatedAt:      currTime,
			UpdatedAt:      currTime,
			JellyfinUserID: u.Id,
			Username:       u.Name,
			IsAdmin:        isAdmin,
		}
		fmt.Println(newUser)
		user, err = cfg.db.AddUser(r.Context(), newUser)
		if err != nil {
			fmt.Println("error creating user: " + err.Error())
			return nil, fmt.Errorf("Error adding new user to db: %v", err.Error())
		}

	} else if err != nil {
		fmt.Println("not sqlerrnorows err: " + err.Error())
		return nil, fmt.Errorf("Error fetching user info: %v", err.Error())
	}

	userVotesCount, err := cfg.db.GetVotesCountPerUser(r.Context(), user.ID)
	if err != nil {
		//TODO: Handle errors check to see if there's an error case for no rows which should just return 0
	}

	userVotedMovies, err := cfg.db.GetMoviesByUserVotes(r.Context(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching movies voted on by the user: %v", err.Error())
	}

	allMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Error fetching all movies for voting: %v", err.Error())
	}

	return templates.VotePage(votedOnMovies, int(userVotesCount), userVotedMovies, allMovies), nil
}

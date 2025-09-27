package main

import (
	"fmt"
	"net/http"
	"strconv"

	"www.github.com/jkboyo/votefin/internal/tmdb"
)

func (cfg *apiConfig) addMovieHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	movieIDStr := r.FormValue("movieID")
	if err != nil {
		//TODO handle error
	}

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing int passed in. \nText passed in%s\nError: %v", movieIDStr, err)
	}
	fmt.Printf("Adding movie with ID %d", movieID)

	movieData, err := tmdb.FetchMovieInfo(movieID)
	if err != nil {
		// TODO: Handle fetch movie info error
	}

	movie, err := cfg.db.InsertMovie(r.Context(), movieData)
	if err != nil {
		// TODO: Handle db insertion error.
	}

}

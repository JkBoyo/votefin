package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"www.github.com/jkboyo/votefin/internal/tmdb"
)

func (cfg *apiConfig) searchMoviesToAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}
	moviePrefix := r.FormValue("moviePrefix")

	tmdbTrie, err := tmdb.InitTMDBTrie()

	currentMatches, err := tmdbTrie.RetrieveObjs(moviePrefix)

}

func (cfg *apiConfig) addMovieHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}
	movieIDStr := r.FormValue("movieID")

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing int passed in. \nText passed in%s\nError: %v", movieIDStr, err)
	}
	fmt.Printf("Adding movie with ID %d\n", movieID)

	movieData, err := tmdb.FetchMovieInfo(movieID)
	if err != nil {
		log.Fatal(err)
	}

	movie, err := cfg.db.InsertMovie(r.Context(), movieData)
	if err != nil {
		log.Fatal(err)
	}

	response := fmt.Sprintf("Successfully added %s for voting", movie.Title)

	respondWithHTMLNotif(w, r, 200, response)
}

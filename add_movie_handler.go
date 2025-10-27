package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"www.github.com/jkboyo/votefin/internal/tmdb"
	"www.github.com/jkboyo/votefin/templates"
)

func (cfg *apiConfig) searchMoviesToAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing forms")
	}

	moviePrefix := r.FormValue("searchMovies")

	tmdbTrie := cfg.tmdbTrie

	if err != nil {
		fmt.Println(err.Error())
		respondWithHtmlErr(w, 500, err.Error())
	}

	fmt.Println(moviePrefix)

	currentMatches, err := tmdbTrie.RetrieveObjs(moviePrefix)
	if err != nil {
		fmt.Println(err.Error())
		respondWithHtmlErr(w, 500, err.Error())
	}

	numRet := min(10, len(currentMatches))

	fmt.Println(numRet)

	fmt.Println(currentMatches[0:numRet])
	retMovies := templates.SearchList(currentMatches[0:numRet])

	fmt.Println(retMovies)

	respondWithHTML(w, 200, retMovies)
}

func (cfg *apiConfig) addMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movieId")

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

	resStr := fmt.Sprintf("Successfully added %s for voting", movie.Title)

	html := templates.AddMovieRet(resStr, movie)

	fmt.Println(html)

	respondWithHTML(w, 200, html)
}

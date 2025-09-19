package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/templates"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	db, err := sql.Open("sqlite3", "./votefin.db")
	if err != nil {
		fmt.Println("Error with DB", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiConf := apiConfig{db: dbQueries}

	userMovies, err := apiConf.db.GetMoviesByUserVotes(context.Background(), 1)
	if err != nil {
		fmt.Println("couldn't fetch movies users voted on movies", err)
	}
	movies, err := apiConf.db.GetMovies(context.Background())
	if err != nil {
		fmt.Println("couldn't fetch movies", err)
	}
	serveMux := http.NewServeMux()

	assets := http.FileServer(http.Dir("./assets/static/"))

	serveMux.Handle("/static/", http.StripPrefix("/static/", assets))
	serveMux.Handle("/", templ.Handler(templates.PageTemplate(true, userMovies, movies)))
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {

}

// I don't think that i will need this but if I change my mind it's here.
// func (cfg *apiConfig) fetchMovies(w http.ResponseWriter, r *http.Request) {
// 	movies, err := cfg.db.GetMovies(r.Context())
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "error fetching movies")
// 		return
// 	}
// 	respondWithJson(w, http.StatusAccepted, movies)
// }
//
// func respondWithJson(w http.ResponseWriter, code int, payload any) error {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		return err
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// 	return nil
// }
//
// func respondWithError(w http.ResponseWriter, code int, msg string) error {
// 	return respondWithJson(w, code, map[string]string{"error": msg})
// }

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/internal/tmdb"
	"www.github.com/jkboyo/votefin/internal/trie"
	"www.github.com/jkboyo/votefin/templates"
)

type apiConfig struct {
	db       *database.Queries
	tmdbTrie *trie.Trie
}

func main() {
	fmt.Println("fetching env data")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loading")
	}
	fmt.Println("loading trie")
	tmdbTrie, err := tmdb.InitTMDBTrie()
	if err != nil {
		log.Fatal("Trie not generating")
	}

	fmt.Println("opening db")
	db, err := sql.Open("sqlite3", "./votefin.db")
	if err != nil {
		fmt.Println("Error with DB", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiConf := apiConfig{db: dbQueries, tmdbTrie: tmdbTrie}

	fmt.Println("fetching movies")
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

	serveMux.HandleFunc("POST /searchmovies", apiConf.searchMoviesToAdd)
	serveMux.HandleFunc("POST /addmovie/", apiConf.addMovieHandler)
	serveMux.Handle("/static/", http.StripPrefix("/static/", assets))
	serveMux.Handle("/", templ.Handler(templates.PageTemplate(true, userMovies, movies)))
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Println("Serving at http://localhost:8080")
	server.ListenAndServe()
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {

}

func respondWithHTML(w http.ResponseWriter, code int, comp templ.Component) error {
	w.Header().Set("Content-Type", "application-/x-www-form-urlencoded")
	w.WriteHeader(code)
	err := comp.Render(context.Background(), os.Stdout)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
	}
	err = comp.Render(context.Background(), w)
	if err != nil {
		respondWithHtmlErr(w, 500, err.Error())
	}
	return nil
}

func respondWithHtmlErr(w http.ResponseWriter, code int, errMsg string) error {
	errNotif := templates.Notification(errMsg)
	return respondWithHTML(w, code, errNotif)
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

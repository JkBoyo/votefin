package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"www.github.com/jkboyo/votefin/internal/database"
	"www.github.com/jkboyo/votefin/internal/tmdb"
	"www.github.com/jkboyo/votefin/internal/trie"
	"www.github.com/jkboyo/votefin/templates"
)

type apiConfig struct {
	db        *database.Queries
	tmdbTrie  *trie.Trie
	voteLimit int
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

	voteLimit, err := strconv.Atoi(os.Getenv("VOTE_LIMIT"))
	if err != nil {
		fmt.Println("No vote limit set defaulting to 5")
		voteLimit = 5
	}

	apiConf := apiConfig{db: dbQueries, tmdbTrie: tmdbTrie, voteLimit: voteLimit}

	serveMux := http.NewServeMux()

	imagesFP := "./assets/static/images/"
	if _, err := os.Stat(imagesFP); os.IsNotExist(err) {
		os.Mkdir(imagesFP, 0775)
	}

	assets := http.FileServer(http.Dir("./assets/static/"))

	serveMux.HandleFunc("POST /searchmovies", apiConf.searchMoviesToAdd)
	serveMux.HandleFunc("POST /addmovie", apiConf.addMovieHandler)
	serveMux.HandleFunc("POST /login", apiConf.loginUser)
	serveMux.HandleFunc("GET /dashboard", apiConf.AuthorizeHandler(apiConf.renderPageHandler))
	serveMux.HandleFunc("POST /vote", apiConf.AuthorizeHandler(apiConf.voteHandler))

	serveMux.Handle("/static/", http.StripPrefix("/static/", assets))

	serveMux.HandleFunc("/{$}", apiConf.AuthorizeHandler(func(w http.ResponseWriter, r *http.Request, user *database.User) {
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusAccepted)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusAccepted)
	}))

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Println("Serving at http://localhost:8080")
	server.ListenAndServe()
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

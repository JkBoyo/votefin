package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	isDev := os.Getenv("DEV_MODE")

	if isDev == "" || isDev == "DEV" {
		fmt.Println("fetching env data")
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(".env not loading")
		}
	}
	err := checkEnvVars()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("loading trie")
	tmdbTrie, err := tmdb.InitTMDBTrie()
	if err != nil {
		slog.Error("error while making trie", "error", err)
		log.Fatal("Trie not generating")
	}

	fmt.Println("opening db")
	db, err := sql.Open("sqlite3", "./data/votefin.db")
	if err != nil {
		fmt.Println("Error with DB", err)
	}
	pragQuery := `PRAGMA journal_mode = WAL;
			PRAGMA synchronous = NORMAL;
			PRAGMA cache_size = 10000;
			PRAGMA temp_store = MEMORY;
			PRAGMA foreign_keys = ON;
			PRAGMA mmap_size = 268435456;`
	_, err = db.Exec(pragQuery)
	if err != nil {
		slog.Error("couldn't execute pragma commands", "error", err)
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

	imagesFP := "./assets/images/"
	if _, err := os.Stat(imagesFP); os.IsNotExist(err) {
		os.Mkdir(imagesFP, 0775)
	}

	serveMux.Handle("/admin/", http.StripPrefix("/admin", apiConf.AuthorizeMiddleWare(CheckIsAdmin(serveMux))))

	serveMux.HandleFunc("POST /searchMovies", apiConf.searchMoviesToAdd)
	serveMux.HandleFunc("POST /addMovie", apiConf.addMovieHandler)
	serveMux.HandleFunc("POST /markFinished", apiConf.markFinishedhandler)
	serveMux.HandleFunc("POST /removeMovie", apiConf.removeMovie)

	serveMux.Handle("POST /login", http.HandlerFunc(apiConf.loginUser))
	serveMux.Handle("POST /logout", http.HandlerFunc(apiConf.logoutUser))

	serveMux.Handle("POST /vote", apiConf.AuthorizeMiddleWare(http.HandlerFunc(apiConf.voteHandler)))
	serveMux.Handle("POST /removeVote", apiConf.AuthorizeMiddleWare(http.HandlerFunc(apiConf.voteRemovalHandler)))

	serveMux.Handle("GET /login", *templ.Handler(templates.BasePage(templates.Login(nil))))
	serveMux.Handle("GET /dashboard", apiConf.AuthorizeMiddleWare(http.HandlerFunc(apiConf.renderPageHandler)))

	serveMux.Handle("/assets/", disableCacheInDevMode(
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))))

	// Hande the root path so that it redirects based on the users cookie
	serveMux.Handle("/{$}", apiConf.AuthorizeMiddleWare(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*database.User)
		if !ok {
			slog.Info("Couldn't find user redirecting to login page")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			fmt.Println("going to login")
			return
		}
		fmt.Println(user)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		fmt.Println("Going to dashboard")
	})))

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Println("Serving at http://localhost:8080")
	server.ListenAndServe()
}

var dev bool = true

func disableCacheInDevMode(next http.Handler) http.Handler {
	if !dev {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func checkEnvVars() error {
	envVars := []string{"VOTE_LIMIT", "SECURE_ONLY_COOKIE", "SERVER_ID", "JELLYFIN_URL", "TMDB_API_KEY", "POSTER_IM_WIDTH", "POPULARITY_LIMIT"}
	varMissing := false
	for _, envVarTitle := range envVars {
		envVar := os.Getenv(envVarTitle)
		if envVar == "" {
			slog.Error("missing variable", "name", envVarTitle)
			varMissing = true
		}
	}
	if varMissing == true {
		return errors.New("Env not complete")
	}

	return nil
}

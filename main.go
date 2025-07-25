package main

import (
	"net/http"

	"www.github.com/jkboyo/votefin/internal/database"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	serveMux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {

}

func (cfg *apiConfig) fetchMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching movies")
		return
	}

}

func respondWithJson(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJson(w, code, map[string]string{"error": msg})
}

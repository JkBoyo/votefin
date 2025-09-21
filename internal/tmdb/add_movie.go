package tmdb

import "www.github.com/jkboyo/votefin/internal/database"


var TMDBBaseURL = "themoviedb.org/"
type tmdbData struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	PosterPath []string `json:"images"`
}

func fetchMovieInfo(tmdbID int64) (database.InsertMovieParams, error) {
	retData := database.InsertMovieParams{
		ID: tmdbID,
		Cr
	}

}

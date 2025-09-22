package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"www.github.com/jkboyo/votefin/internal/database"
)

var authHeader = fmt.Sprintf("Bearer %s", os.Getenv("TMDB_API_KEY"))

var TMDBBaseURL = "themoviedb.org/"
var TMDBImageBaseUrl = "https://imgage.themoviedb.org/t/p/w342%s"

type tmdbData struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	PosterPath string `json:"images"`
}

func fetchMovieInfo(tmdbID int64) (database.InsertMovieParams, error) {
	now := time.Now().Local().GoString()
	infoReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.%s3/movie/%d", TMDBBaseURL, tmdbID), nil)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Movie request failed with error: %v", err)
	}
	infoReq.Header.Add("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(infoReq)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Movie response failed with error: %v", err)
	}

	defer resp.Body.Close()

	var data *tmdbData

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Error reading response bytes: %v", err)
	}

	if err := json.Unmarshal(respData, &data); err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Error unmarshaling json: %v", err)
	}

	go writeImageToDisk(data.PosterPath)

	retData := database.InsertMovieParams{
		ID:         tmdbID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Title:      data.Title,
		TmdbUrl:    fmt.Sprintf("https://www.%s/%d", TMDBBaseURL, data.ID),
		PosterPath: data.PosterPath,
		Status:     "Not ripped.",
	}

	return retData, nil
}

func writeImageToDisk(pp string) {
	// TODO Better error formatting
	imageUrl := fmt.Sprintf(TMDBImageBaseUrl, pp)

	resp, err := http.DefaultClient.Get(imageUrl)
	if err != nil {
		log.Print(err)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	err = os.WriteFile(pp, dat, 0775)
	if err != nil {
		log.Print(err)
	}
}

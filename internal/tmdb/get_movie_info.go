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

var TMDBBaseURL = "themoviedb.org/"
var TMDBImageBaseUrl = "https://imgage.themoviedb.org/t/p/w342%s"

type tmdbData struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	PosterPath string `json:"images"`
}

func FetchMovieInfo(tmdbID int64) (database.InsertMovieParams, error) {
	authHeader := fmt.Sprintf("Bearer %s", os.Getenv("TMDB_API_KEY"))
	now := time.Now().Local().GoString()
	movieURL := fmt.Sprintf("https://api.%s3/movie/%d", TMDBBaseURL, tmdbID)
	infoReq, err := http.NewRequest("GET", movieURL, nil)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Movie request failed with error: %v\n", err)
	}
	fmt.Println(authHeader)
	infoReq.Header.Add("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(infoReq)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Movie response failed with error: %v\n", err)
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()

	var data *tmdbData

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Error reading response bytes: %v\n", err)
	}

	if err := json.Unmarshal(respData, &data); err != nil {
		return database.InsertMovieParams{}, fmt.Errorf("Error unmarshaling json: %v\n", err)
	}

	fmt.Println(*data)

	retData := database.InsertMovieParams{
		ID:         tmdbID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Title:      data.Title,
		TmdbUrl:    fmt.Sprintf("https://www.%s/%d", TMDBBaseURL, data.ID),
		PosterPath: data.PosterPath,
		Status:     "Not ripped.",
	}

	fmt.Println(retData.PosterPath)
	writeImageToDisk(retData.PosterPath)

	return retData, nil
}

func writeImageToDisk(pp string) {
	if pp == "" {
		log.Fatal("No Poster Path provided")
	}
	// TODO: Better error formatting
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

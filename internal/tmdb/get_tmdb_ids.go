package tmdb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetTMDBDataFP() (string, error) {
	today := time.Now()
	dataDirEntries, err := os.ReadDir("./data/")
	if err != nil {
		return "", err
	}

	currentMovieDataName := ""
	currentMovieDataDate := time.Time{}
	movieIdFileCount := 0

	for _, dirEntry := range dataDirEntries {
		if dirEntry.IsDir() || filepath.Ext(dirEntry.Name()) != ".json" {
			continue
		}
		movieIdFileCount += 1

		fileInfo, err := dirEntry.Info()
		if err != nil {
			return "", err
		}
		fileDate, err := getTMDBIDFileDate(fileInfo)
		if err != nil {
			return "", err
		}
		if fileDate.After(currentMovieDataDate) {
			currentMovieDataName = dirEntry.Name()
			currentMovieDataDate = fileDate
		}
	}
	if currentMovieDataDate.Before(today.Add(-time.Hour*24*7)) || movieIdFileCount == 0 {
		err = downloadCurrentTMDBIDFile(today)
		if err != nil {
			return "", err
		}
	}

	tmdbDataFP := "./data/" + currentMovieDataName

	return tmdbDataFP, nil

}

func downloadCurrentTMDBIDFile(today time.Time) error {
	yesterday := today.Add(time.Hour * 24)
	currIdFileName := fmt.Sprintf("movie_ids_%d_%d_%d.json.gz",
		int(yesterday.Month()),
		yesterday.Day(),
		yesterday.Year(),
	)
	currIdUrl := "https://files.tmdb.org/p/exports/" + currIdFileName
	resp, err := http.DefaultClient.Get(currIdUrl)
	if err != nil {
		return err
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	idFileFP := "./data/" + currIdFileName

	err = os.WriteFile(idFileFP, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getTMDBIDFileDate(fileInfo os.FileInfo) (time.Time, error) {
	prefixLessString := strings.TrimPrefix(fileInfo.Name(), "movie_ids_")
	trimmedString := strings.TrimSuffix(prefixLessString, ".json")

	fileDate, err := time.Parse(
		trimmedString,
		"01_02_2006",
	)
	return fileDate, err
}

package tmdb

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
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
	if currentMovieDataDate.Before(today.AddDate(0, 0, -7)) || movieIdFileCount == 0 {
		err = downloadCurrentTMDBIDFile(today)
		if err != nil {
			return "", err
		}
	}

	tmdbDataFP := "./data/" + currentMovieDataName

	return tmdbDataFP, nil

}

func downloadCurrentTMDBIDFile(today time.Time) error {
	yesterday := today.AddDate(0, 0, -1)
	currIdFileName := fmt.Sprintf("movie_ids_%02d_%d_%d.json",
		int(yesterday.Month()),
		yesterday.Day(),
		yesterday.Year(),
	)
	currIdUrl := "https://files.tmdb.org/p/exports/" + currIdFileName + ".gz"
	resp, err := http.DefaultClient.Get(currIdUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		slog.Info("Request failed", "url", resp.Request.URL)
		return fmt.Errorf("Response failed with status: %s", resp.Status)
	}

	idFileFP := "./data/" + currIdFileName
	idFileFPgz := idFileFP + ".gz"

	tempFile, err := os.Create(idFileFPgz)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	gzFile, err := os.Open(idFileFPgz)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	reader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	finalFile, err := os.Create(idFileFP)
	if err != nil {
		return err
	}

	_, err = io.Copy(finalFile, reader)
	if err != nil {
		return err
	}
	err = os.Remove(idFileFPgz)
	if err != nil {
		return err
	}
	return nil
}

func getTMDBIDFileDate(fileInfo os.FileInfo) (time.Time, error) {
	prefixLessString := strings.TrimPrefix(fileInfo.Name(), "movie_ids_")
	trimmedString := strings.TrimSuffix(prefixLessString, ".json")

	fileDate, err := time.Parse(
		"01_02_2006",
		trimmedString,
	)
	return fileDate, err
}

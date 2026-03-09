package tmdb

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"www.github.com/jkboyo/votefin/internal/trie"
)

func InitTMDBTrie() (*trie.Trie, error) {
	tmdbTrie := trie.NewTrie()

	tmdbDataFP := "./data/movie_ids_08_30_2025.json"

	tmdbData, err := os.Open(tmdbDataFP)

	if err != nil {
		return nil, errors.New("error opening TMDB DATA file")
	}

	defer tmdbData.Close()

	scanner := bufio.NewScanner(tmdbData)

	scanner.Split(bufio.ScanLines)

	popLimit, err := strconv.ParseFloat(os.Getenv("POPULARITY_LIMIT"), 64)
	if err != nil {
		popLimit = 0
	}

	for scanner.Scan() {
		newMovie := &trie.Movie{}
		err := json.Unmarshal(scanner.Bytes(), newMovie)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling: %s", err.Error())
		}
		if newMovie.Popularity < float32(popLimit) {
			continue
		}
		tmdbTrie.Insert(*newMovie)
	}

	return tmdbTrie, nil
}

// func GetMovieTrie() error {}

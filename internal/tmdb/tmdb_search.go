package tmdb

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"www.github.com/jkboyo/votefin/internal/trie"

	"github.com/joho/godotenv"
)

func InitTMDBTrie() (*trie.Trie, error) {
	tmdbTrie := trie.NewTrie()

	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("TMDB data env var not set.")
	}
	tmdbDataFP := os.Getenv("TMDB_DATA")

	tmdbData, err := os.Open(tmdbDataFP)

	if err != nil {
		return nil, errors.New("error opening TMDB DATA file")
	}

	defer tmdbData.Close()

	scanner := bufio.NewScanner(tmdbData)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		newObj := &struct {
			ID            int     `json:"id"`
			OriginalTitle string  `json:"original_title"`
			Popularity    float32 `json:"popularity"`
		}{}
		err := json.Unmarshal(scanner.Bytes(), newObj)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling: %s", err.Error())
		}
		trieObj := trie.Obj{
			Str:        newObj.OriginalTitle,
			Val:        newObj.ID,
			Popularity: newObj.Popularity,
		}
		tmdbTrie.Insert(trieObj)
	}

	return tmdbTrie, nil
}

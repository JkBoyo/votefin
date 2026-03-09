package trie

import (
	"errors"
	"maps"
	"slices"
	"strings"
)

var (
	ErrNoMatch = errors.New("No objects in trie match")
)

type trieNode struct {
	Movies   []*Movie
	Children map[rune]*trieNode
}

type Trie struct {
	Root *trieNode
}

type Obj struct {
	Str        string
	Val        int
	Popularity float32
}

type Movie struct {
	Title      string  `json:"original_title"`
	ID         int     `json:"id"`
	Popularity float32 `json:"popularity"`
}

func NewTrie() *Trie {
	return &Trie{
		Root: &trieNode{
			Children: map[rune]*trieNode{},
		},
	}
}

func (t *Trie) Insert(movie Movie) {
	curNode := t.Root
	simpObjStr := reduceTrieStr(movie.Title)
	for i, char := range simpObjStr {
		nextNode, exists := curNode.Children[char]
		if exists {
			curNode = nextNode
			continue
		}

		newNode := &trieNode{}
		if i < len(simpObjStr)-1 {
			newNode.Children = map[rune]*trieNode{}
			curNode.Children[char] = newNode
			curNode = newNode
		}

	}
	curNode.Movies = append(curNode.Movies, &movie)
}

func (t *Trie) RetrieveObjs(pref string) ([]*Movie, error) {
	curNode := t.Root
	simpPref := reduceTrieStr(pref)
	for _, char := range simpPref {
		nextNode, exists := curNode.Children[char]
		if !exists {
			return []*Movie{}, ErrNoMatch
		}
		curNode = nextNode
	}

	retObjs := searchLevel(curNode, pref)
	if curNode.Movies != nil {
		for _, movie := range curNode.Movies {
			retObjs = append(retObjs, movie)
		}
	}

	slices.SortFunc(retObjs, func(i, j *Movie) int {
		return int(j.Popularity) - int(i.Popularity)
	})

	return retObjs, nil
}

func searchLevel(currNode *trieNode, currPrefix string) []*Movie {
	keys := maps.Keys(currNode.Children)

	movies := []*Movie{}

	for k := range keys {
		if currNode.Children[k].Movies != nil {
			for _, movie := range currNode.Children[k].Movies {
				movies = append(movies, movie)
			}
		}
		if currNode.Children[k].Children != nil {
			movies = append(movies, searchLevel(currNode.Children[k], currPrefix+string(k))...)
		}
	}
	return movies
}

func reduceTrieStr(str string) string {
	lowerStr := strings.ToLower(str)
	spaceLessStr := strings.ReplaceAll(lowerStr, " ", "")
	stopWordLessStr := spaceLessStr
	// Remove all stop words defined in this list
	// Stop words need to use regex to make sure things are correctly
	stopWords := []string{"the", "and", "of", ":"}
	for _, word := range stopWords {
		stopWordLessStr = strings.ReplaceAll(stopWordLessStr, word, "")
	}

	return stopWordLessStr
}

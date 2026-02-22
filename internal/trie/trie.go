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
	Children  map[rune]*trieNode
	Movies    []Movie
	IsNameEnd bool
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
	Title      string
	ID         int
	Popularity float32
}

func NewTrie() *Trie {
	return &Trie{
		Root: &trieNode{
			Children:  map[rune]*trieNode{},
			Movies:    []Movie{},
			IsNameEnd: false,
		},
	}
}

func (t *Trie) Insert(obj Obj) {
	curNode := t.Root
	simpObjStr := reduceTrieStr(obj.Str)
	for _, char := range simpObjStr {
		nextNode, exists := curNode.Children[char]
		if exists {
			curNode = nextNode
		} else {
			newNode := &trieNode{
				Children:  map[rune]*trieNode{},
				IsNameEnd: false,
			}
			curNode.Children[char] = newNode
			curNode = newNode
		}
	}
	curNode.Movies = append(curNode.Movies, Movie{Title: obj.Str, ID: obj.Val, Popularity: obj.Popularity})
	curNode.IsNameEnd = true
}

func (t *Trie) RetrieveObjs(pref string) ([]Obj, error) {
	curNode := t.Root
	simpPref := reduceTrieStr(pref)
	for _, char := range simpPref {
		nextNode, exists := curNode.Children[char]
		if !exists {
			return []Obj{}, ErrNoMatch
		}
		curNode = nextNode
	}

	retObjs := searchLevel(curNode, pref)
	if curNode.IsNameEnd {
		for _, movie := range curNode.Movies {
			retObjs = append(retObjs, Obj{movie.Title, movie.ID, movie.Popularity})
		}
	}

	slices.SortFunc(retObjs, func(i, j Obj) int {
		return int(j.Popularity) - int(i.Popularity)
	})

	return retObjs, nil
}

func searchLevel(currNode *trieNode, currPrefix string) []Obj {
	keys := maps.Keys(currNode.Children)

	objs := []Obj{}

	for k := range keys {
		if currNode.Children[k].IsNameEnd {
			for _, movie := range currNode.Children[k].Movies {
				objs = append(objs, Obj{movie.Title, movie.ID, movie.Popularity})
			}
		}
		if currNode.Children[k].Children != nil {
			objs = append(objs, searchLevel(currNode.Children[k], currPrefix+string(k))...)
		}
	}
	return objs
}

func reduceTrieStr(str string) string {
	lowerStr := strings.ToLower(str)
	spaceLessStr := strings.ReplaceAll(lowerStr, " ", "")
	stopWordLessStr := spaceLessStr
	// Remove all stop words defined in this list
	stopWords := []string{"the", "and", "of", ":"}
	for _, word := range stopWords {
		stopWordLessStr = strings.ReplaceAll(stopWordLessStr, word, "")
	}

	return stopWordLessStr
}

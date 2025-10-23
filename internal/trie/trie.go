package trie

import (
	"errors"
	"maps"
	"slices"
)

var (
	ErrNoMatch = errors.New("No objects in trie match")
)

type trieNode struct {
	Children   map[rune]*trieNode
	Id         int
	Popularity float32
	IsNameEnd  bool
}

type Trie struct {
	Root *trieNode
}

type Obj struct {
	Str        string
	Val        int
	Popularity float32
}

func NewTrie() *Trie {
	return &Trie{
		Root: &trieNode{
			Children:  map[rune]*trieNode{},
			Id:        0,
			IsNameEnd: false,
		},
	}
}

func (t *Trie) Insert(obj Obj) {
	curNode := t.Root
	for _, char := range obj.Str {
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
	curNode.Id = obj.Val
	curNode.Popularity = obj.Popularity
	curNode.IsNameEnd = true
}

func (t *Trie) RetrieveObjs(pref string) ([]Obj, error) {
	curNode := t.Root
	for _, char := range pref {
		nextNode, exists := curNode.Children[char]
		if !exists {
			return []Obj{}, ErrNoMatch
		}
		curNode = nextNode
	}

	retObjs := searchLevel(curNode, pref)

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
			objs = append(objs, Obj{currPrefix + string(k), currNode.Children[k].Id, currNode.Children[k].Popularity})
		}
		if currNode.Children[k].Children != nil {
			objs = append(objs, searchLevel(currNode.Children[k], currPrefix+string(k))...)
		}
	}
	return objs
}

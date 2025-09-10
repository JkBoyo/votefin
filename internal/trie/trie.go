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
	remChil := curNode.Children

	retObjs := searchLevel(remChil, pref)
	slices.SortFunc(retObjs, func(i, j Obj) int {
		if i.Popularity > j.Popularity {
			return -1
		}
		if i.Popularity < j.Popularity {
			return 1
		}
		if i.Popularity == j.Popularity {
			return 0
		}
		return 0
	})
	return retObjs, nil
}

func searchLevel(currLev map[rune]*trieNode, currPrefix string) []Obj {
	keys := maps.Keys(currLev)

	objs := []Obj{}

	for k := range keys {
		if currLev[k].IsNameEnd {
			objs = append(objs, Obj{currPrefix + string(k), currLev[k].Id, currLev[k].Popularity})
		}
		if len(currLev[k].Children) != 0 {
			objs = append(objs, searchLevel(currLev[k].Children, currPrefix+string(k))...)
		}
	}
	return objs
}

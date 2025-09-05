package trie

import (
	"errors"
	"maps"
)

var (
	ErrNoMatch = errors.New("No objects in trie match")
)

type trieNode struct {
	Children  map[rune]*trieNode
	Id        int
	IsNameEnd bool
}

type Trie struct {
	Root *trieNode
}

type Obj struct {
	Str string
	Val int
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

func (t *Trie) Insert(obj string, id int) {
	curNode := t.Root
	for _, char := range obj {
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
	curNode.Id = id
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

	return searchLevel(remChil, pref), nil
}

func searchLevel(currLev map[rune]*trieNode, currPrefix string) []Obj {
	keys := maps.Keys(currLev)

	objs := []Obj{}

	for k := range keys {
		if currLev[k].IsNameEnd {
			objs = append(objs, Obj{currPrefix + string(k), currLev[k].Id})
		}
		if len(currLev[k].Children) != 0 {
			objs = append(objs, searchLevel(currLev[k].Children, currPrefix+string(k))...)
		}
	}
	return objs
}

package trie

import (
	"errors"
	"maps"
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
			newNode := trieNode{
				Children:  map[rune]*trieNode{},
				IsNameEnd: false,
			}
			curNode.Children[char] = &newNode
			curNode = &newNode
		}
	}
	newNode := trieNode{
		Children:  nil,
		Id:        id,
		IsNameEnd: true,
	}
	curNode.Children[rune(0)] = &newNode
}

func (t *Trie) RetrieveObjs(obj string, numRet int) ([]Obj, error) {
	curNode := t.Root
	for _, char := range obj {
		nextNode, exists := curNode.Children[char]
		if !exists {
			return []Obj{}, errors.New("No objects in trie match")
		}
		curNode = nextNode
	}
	remChil := curNode.Children

	objs := []Obj{}
	return searchLevel(remChil, obj, objs, numRet), nil
}

func searchLevel(currLev map[rune]*trieNode, currPrefix string, objs []Obj, numRet int) []Obj {
	keys := maps.Keys(currLev)
	for r := range keys {
		if len(objs) == numRet {
			return objs
		}
		if r == rune(0) {
			objs = append(objs, Obj{Str: currPrefix, Val: currLev[r].Id})
			continue
		}
		objs = append(objs, searchLevel(currLev[r].Children, currPrefix+string(r), objs, numRet)...)
	}
	return objs
}

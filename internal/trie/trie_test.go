package trie

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestTrie(t *testing.T) {
	testTrie := NewTrie()
	testObjSlice := []Obj{
		{Str: "This word is cool", Val: 123456, Popularity: 10.0},
		{Str: "This sentence is different", Val: 54321, Popularity: 7.0},
		{Str: "This sentence has some diff", Val: 56789, Popularity: 5.0},
		{Str: "This sentence has some differences", Val: 78329, Popularity: 6.0},
		{Str: "This sentence has some differences", Val: 3456789, Popularity: 4.0},
		{Str: "This word drools", Val: 98765, Popularity: 3.0},
		{Str: "This word is happy", Val: 56466, Popularity: 1.0},
	}
	for _, obj := range testObjSlice {
		testTrie.Insert(obj)
	}
	printTrie, err := json.MarshalIndent(testTrie, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(printTrie))

	testCases := []struct {
		expected []Obj
		err      error
		name     string
		prefix   string
	}{
		{
			[]Obj{
				{Str: "This word is cool", Val: 123456, Popularity: 10.0},
				{Str: "This sentence is different", Val: 54321, Popularity: 7.0},
				{Str: "This sentence has some differences", Val: 78329, Popularity: 6.0},
				{Str: "This sentence has some diff", Val: 56789, Popularity: 5.0},
				{Str: "This sentence has some differences", Val: 3456789, Popularity: 4.0},
				{Str: "This word drools", Val: 98765, Popularity: 3.0},
				{Str: "This word is happy", Val: 56466, Popularity: 1.0},
			},
			nil,
			"All Ojbs should exist and be returned",
			"This",
		},
		{
			[]Obj{
				{Str: "This sentence is different", Val: 54321, Popularity: 7.0},
				{Str: "This sentence has some differences", Val: 78329, Popularity: 6.0},
				{Str: "This sentence has some diff", Val: 56789, Popularity: 5.0},
				{Str: "This sentence has some differences", Val: 3456789, Popularity: 4.0},
			},
			nil,
			"Objs with 'This sentence' prefixes should be returned",
			"This sentence",
		},
		{
			[]Obj{},
			ErrNoMatch,
			"No objects in the trie match the prefix",
			"That",
		},
	}

	for _, tC := range testCases {
		results, err := testTrie.RetrieveObjs(tC.prefix)

		fmt.Println(tC.expected)
		fmt.Println(results)
		if !reflect.DeepEqual(results, tC.expected) {
			t.Errorf("\n%s\n failed obj's.\n Expected: %v\n Actual: %v",
				tC.name, tC.expected, results)
		}
		if !errors.Is(err, tC.err) {
			t.Errorf("\n%s\n failed on errors.\n Expected: %v\n Actual: %v",
				tC.name, tC.err, err)
		}
	}
}

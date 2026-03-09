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
	testMovieSlice := []Movie{
		{Title: "This word is cool", ID: 123456, Popularity: 10.0},
		{Title: "This sentence is different", ID: 54321, Popularity: 7.0},
		{Title: "This sentence has some diff", ID: 56789, Popularity: 5.0},
		{Title: "This sentence has some differences", ID: 78329, Popularity: 6.0},
		{Title: "This sentence has some differences", ID: 3456789, Popularity: 4.0},
		{Title: "This word drools", ID: 98765, Popularity: 3.0},
		{Title: "This word is happy", ID: 56466, Popularity: 1.0},
	}
	for _, obj := range testMovieSlice {
		testTrie.Insert(obj)
	}
	printTrie, err := json.MarshalIndent(testTrie, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(printTrie))

	testCases := []struct {
		expected []*Movie
		err      error
		name     string
		prefix   string
	}{
		{
			[]*Movie{
				{Title: "This word is cool", ID: 123456, Popularity: 10.0},
				{Title: "This sentence is different", ID: 54321, Popularity: 7.0},
				{Title: "This sentence has some differences", ID: 78329, Popularity: 6.0},
				{Title: "This sentence has some diff", ID: 56789, Popularity: 5.0},
				{Title: "This sentence has some differences", ID: 3456789, Popularity: 4.0},
				{Title: "This word drools", ID: 98765, Popularity: 3.0},
				{Title: "This word is happy", ID: 56466, Popularity: 1.0},
			},
			nil,
			"All Ojbs should exist and be returned",
			"This",
		},
		{
			[]*Movie{
				{Title: "This sentence is different", ID: 54321, Popularity: 7.0},
				{Title: "This sentence has some differences", ID: 78329, Popularity: 6.0},
				{Title: "This sentence has some diff", ID: 56789, Popularity: 5.0},
				{Title: "This sentence has some differences", ID: 3456789, Popularity: 4.0},
			},
			nil,
			"Movies with 'This sentence' prefixes should be returned",
			"This sentence",
		},
		{
			[]*Movie{},
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

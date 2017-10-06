package main

import (
	"testing"
)

func Test_persistanceTest(t *testing.T) {
	id := uint64(123)
	text := "hello world"
	OpenDB(true)
	Insert(id, text)
	CloseDB()

	OpenDB(true)
	textFromSelect, found := Select(id)
	CloseDB()

	if !found || text != textFromSelect {
		t.Errorf("persistance broken")
	}
	ResetDB()
}

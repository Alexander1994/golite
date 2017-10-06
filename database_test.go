package main

import (
	"os"
	"testing"
)

func setupTests() {
	OpenDB(true)
}

func tearDownTests() {
	CloseDB()
	DeleteDB()
}

func TestMain(m *testing.M) {
	setupTests()
	retCode := m.Run()
	tearDownTests()
	os.Exit(retCode)
}

func Test_persistanceTest(t *testing.T) {
	id := uint64(123)
	text := "hello world"
	isInsert := Insert(id, text)
	CloseDB()
	if !isInsert {
		t.Errorf("could not insert")
	}

	OpenDB(true)
	textFromSelect, found := Select(id)

	if !found || text != textFromSelect {
		t.Errorf("persistance broken")
	}
}

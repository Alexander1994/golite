package main

import (
	"os"
	"testing"
)

var _1000CharString string

func setupTests() {
	OpenDB(true)
	_1000CharString = stringOfLength(1000)
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

func Benchmark1000lengthInsert(b *testing.B) {
	Insert(1, _1000CharString)
}

// helper functions
func stringOfLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

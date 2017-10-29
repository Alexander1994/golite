package main

import (
	"os"
	"testing"
)

var _1000CharString string

func setupTests() {
	OpenDB()
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
	id := uint32(123)
	text := "hello world"
	tableName := "test"
	if !CreateTable(tableName) {
		t.Errorf("could not create a table")
	}
	isInsert := Insert(id, text, tableName)
	CloseDB()
	if !isInsert {
		t.Errorf("could not insert")
	}

	OpenDB()
	textFromSelect, found := Select(id, tableName)

	if !found || text != textFromSelect {
		t.Errorf("persistance broken")
	}

	if !Delete(id, tableName) {
		t.Errorf("Delete broken")
	}
}

func Benchmark1000lengthInsert(b *testing.B) {
	tableName := "test"
	Insert(1, _1000CharString, tableName)
}

// helper functions
func stringOfLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

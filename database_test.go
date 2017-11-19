package main

import (
	"os"
	"testing"

	"github.com/Alexander1994/golite/database"
)

var _1000CharString string

func setupTests() {
	database.OpenDB()
	_1000CharString = stringOfLength(1000)
}

func tearDownTests() {
	database.CloseDB()
	database.DeleteDB()
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
	if !database.CreateTable(tableName) {
		t.Errorf("could not create a table")
	}
	isInsert := database.Insert(id, text, tableName)
	database.CloseDB()
	if !isInsert {
		t.Errorf("could not insert")
	}

	database.OpenDB()
	textFromSelect, found := database.Select(id, tableName)

	if !found || text != textFromSelect {
		t.Errorf("persistance broken")
	}

	if !database.Delete(id, tableName) {
		t.Errorf("Delete broken")
	}
}

func Benchmark_1000lengthInsert(b *testing.B) {
	tableName := "test"
	database.Insert(1, _1000CharString, tableName)
}

// helper functions
func stringOfLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

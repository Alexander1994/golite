package main

import (
	"os"
	"testing"

	"github.com/Alexander1994/golite/database"
)

var _1000CharString string
var tableName string

// setup suite and test
func setupTestSuite() {
	_1000CharString = stringOfLength(1000)
	tableName = "test"
}

func tearDownTestSuite() {

}

func setupTest() {
	database.OpenDB()
	database.DeleteTable(tableName)
	database.CreateTable(tableName)
}

func tearDownTest() {
	database.CloseDB()
	database.DeleteDB()
}

func TestMain(m *testing.M) {
	setupTestSuite()
	retCode := m.Run()
	tearDownTestSuite()
	os.Exit(retCode)
}

// Tests
func Test_persistanceTest(t *testing.T) {
	setupTest()
	id := uint32(123)
	text := "hello world"
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
	tearDownTest()
}

// Benchmarks
func Benchmark_1000lengthInsert(b *testing.B) {
	setupTest()
	database.Insert(1, _1000CharString, tableName)
	tearDownTest()
}

// helper functions
func stringOfLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

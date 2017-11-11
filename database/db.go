package database

import (
	"io/ioutil"
	"os"
)

const dirname = ".data"

type database map[string]table

var db = make(database)

func openDB() {
	os.Mkdir(dirname, 0755)
	files, err := ioutil.ReadDir(".data")
	fatal(err)

	for _, f := range files {
		if !f.IsDir() {
			fileName := f.Name()
			name := fileName[0 : len(fileName)-4]
			db[name] = createTable(name)
		}
	}
}

func closeDB() {
	for _, dbTable := range db {
		dbTable.close()
	}
}

func deleteDB() {
	for _, dbTable := range db {
		dbTable.delete()
	}
}

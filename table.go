package main

import "os"

type table struct {
	name  string
	file  disk
	cache cacheTable
}

func createTable(name string) table {
	file := createDisk(name)
	return table{name, *file, createCache()}
}

func (dbTable table) close() {
	dbTable.closeCache()
	dbTable.cache.reset()
	dbTable.file.close()
}

func (dbTable table) delete() {
	dbTable.cache.reset()
	dbTable.file.close()
	os.Remove(dbTable.getFileName())
}

func (dbTable table) closeCache() {
	for id, cacheRow := range dbTable.cache {
		if !cacheRow.inMem {
			dbTable.file.pushToDisk(id, cacheRow.text)
		}
	}
}

func (dbTable table) getFileName() string {
	return dirname + "/" + dbTable.name + ".dat"
}

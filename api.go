package main

// OpenDB is for opening the DB. found in database.go
func OpenDB(testMode bool) {
	openDB()
}

// CloseDB is for closing the DB. found in database.go
func CloseDB() {
	closeDB()
}

// DeleteDB clears the entire db
func DeleteDB() {
	deleteDB()
}

// CreateTable creates a table
func CreateTable(name string) bool {
	_, found := db[name]
	if !found {
		db[name] = createTable(name)
	}
	return !found
}

// DeleteTable deletes the named table
func DeleteTable(name string) bool {
	table, found := db[name]
	if found {
		table.delete()
		delete(db, name)
	}
	return found
}

// Insert is for inserting a row into the db
func Insert(id uint32, text string, tableName string) bool { // succesful insert
	dbTable, found := db[tableName]
	if !found {
		return false
	}
	cache := dbTable.cache
	_, foundInCache := cache[id]
	if foundInCache { // if ID is in the cache
		print("blu")
		return false
	}
	if len(cache) < maxCacheSize {
		cache.addRow(id, text)
		return true
	}
	lowestCacheHitRateID := cache.getLowestHitRowID()
	rowToPushToMemory := cache[lowestCacheHitRateID]
	cache.deleteRow(lowestCacheHitRateID)
	if !rowToPushToMemory.inMem {
		cache.addRow(id, text)
		return dbTable.file.pushToDisk(lowestCacheHitRateID, rowToPushToMemory.text) // if ID already exists in memory
	}
	return true
}

// Delete is for deleting a row from the db
func Delete(id uint32, tableName string) bool {
	dbTable, found := db[tableName]
	if !found {
		return false
	}
	cache := dbTable.cache
	file := dbTable.file
	cacheRow, textFound := cache[id]

	if textFound {
		delete(cache, id)
		if cacheRow.inMem {
			file.deleteRowFromDisk(id)
		}
		return true
	}
	return file.deleteRowFromDisk(id)
}

// Select is for retrieving a row from the db
func Select(id uint32, tableName string) (string, bool) {
	dbTable, found := db[tableName]
	if !found {
		return "", false
	}
	cache := dbTable.cache
	file := dbTable.file
	cacheRow, textFound := cache[id]
	if textFound {
		cacheRow.selectCount++
		return cacheRow.text, true
	} else if file.size() > 0 {
		text, found := file.getRowFromDisk(id)
		if found {
			if len(cache) == maxCacheSize {
				lowestCacheHitRateID := cache.getLowestHitRowID()
				cache.deleteRow(lowestCacheHitRateID)
				if !cache[lowestCacheHitRateID].inMem {
					file.pushToDisk(lowestCacheHitRateID, cache[lowestCacheHitRateID].text)
				}
			}
			cache.addMemRow(id, text)
			return text, true
		}
	}
	return "", false
}

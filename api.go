package main

import "os"

// OpenDB is for opening the DB. found in database.go
func OpenDB(testMode bool) {
	file.open(testMode)
}

// CloseDB is for closing the DB. found in database.go
func CloseDB() {
	closeCache()
	resetCache()
	file.close()
}

// DeleteDB clears the entire db
func DeleteDB() {
	resetCache()
	file.close()
	os.Remove(fileName)
}

// Insert is for inserting a row into the db
func Insert(id uint32, text string) bool { // succesful insert
	_, foundInCache := cache[id]
	if foundInCache { // if ID is in the cache
		return false
	}
	if len(cache) < maxCacheSize {
		addCacheRow(id, text)
		return true
	}
	lowestCacheHitRateID := getLowestHitRowID()
	rowToPushToMemory := cache[lowestCacheHitRateID]
	deleteRowFromCache(lowestCacheHitRateID)
	if !rowToPushToMemory.inMem {
		addCacheRow(id, text)
		return file.pushToDisk(lowestCacheHitRateID, rowToPushToMemory.text) // if ID already exists in memory
	}
	return true
}

// Delete is for deleting a row from the db
func Delete(id uint32) bool {
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
func Select(id uint32) (string, bool) {
	cacheRow, textFound := cache[id]
	if textFound {
		cacheRow.selectCount++
		return cacheRow.text, true
	} else if file.size() > 0 {
		text, found := file.getRowFromDisk(id)
		if found {
			if len(cache) == maxCacheSize {
				lowestCacheHitRateID := getLowestHitRowID()
				deleteRowFromCache(lowestCacheHitRateID)
				if !cache[lowestCacheHitRateID].inMem {
					file.pushToDisk(lowestCacheHitRateID, cache[lowestCacheHitRateID].text)
				}
			}
			addMemRowToCache(id, text)
			return text, true
		}
	}
	return "", false
}

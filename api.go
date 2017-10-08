package main

import "os"

// OpenDB is for opening the DB. found in database.go
func OpenDB(testMode bool) {
	openDisk(testMode)
	loadPageTable()
}

// CloseDB is for closing the DB. found in database.go
func CloseDB() {
	closeCache()
	resetCache()
	resetPageTable()
	closeDisk()
}

// DeleteDB clears the entire db
func DeleteDB() {
	resetCache()
	resetPageTable()
	closeDisk()
	os.Remove(fileName)
}

// Insert is for inserting a row into the db
func Insert(id uint64, text string) bool {
	_, foundInCache := cache[id]
	_, foundOnDisk := getRowFromDisk(id)
	inDB := foundInCache || foundOnDisk
	if !inDB {
		if len(cache) < maxCacheSize {
			addCacheRow(id, text)
		} else {
			lowestCacheHitRateID := getLowestHitRowID()
			rowToPushToMemory := cache[lowestCacheHitRateID]
			if !rowToPushToMemory.inMem {
				addCacheRow(id, text)
				pushToDisk(lowestCacheHitRateID, rowToPushToMemory.text)
			}
			deleteRowFromCache(lowestCacheHitRateID)
		}
	}
	return !inDB // if not in DB then it was inserted
}

// Delete is for deleting a row from the db
func Delete(id uint64) bool {
	cacheRow, textFound := cache[id]
	if textFound {
		delete(cache, id)
		if cacheRow.inMem {
			deleteRowFromDisk(id)
		}
		return true
	}
	return deleteRowFromDisk(id)
}

// Select is for retrieving a row from the db
func Select(id uint64) (string, bool) {
	cacheRow, textFound := cache[id]
	fileStat, _ := file.Stat()
	if textFound {
		cacheRow.selectCount++
		return cacheRow.text, true
	} else if fileStat.Size() > 0 {
		memoryRow, found := getRowFromDisk(id)
		if found {
			if len(cache) == maxCacheSize {
				lowestCacheHitRateID := getLowestHitRowID()
				deleteRowFromCache(lowestCacheHitRateID)
			}
			addMemoryRowToCacheTable(memoryRow)
			return memoryRow.text, true
		}
	}
	return "", false
}

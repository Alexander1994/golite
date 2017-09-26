
package main

import (
	"fmt"
)

const MAX_ROW_BYTES = 1024
const MAX_CACHE_SIZE = 20

const MAX_UINT64 = ^uint64(0)

type CacheRow struct {
	text string
	selectCount uint64
	inMem bool
}

var cache = make(map[uint64]CacheRow)


// Insert Into Cache
func insertRow(id uint64, text []string) {
	textToAddToRow := ""
	for i, param := range text {
		if i != len(text) - 1 {
			param += " "
		}
		textToAddToRow += param
	}
	if len(cache) < MAX_CACHE_SIZE {
		addCacheRow(id, textToAddToRow)
	} else {
		lowestCacheHitRateId := getLowestHitRowId()
		rowToPushToMemory := cache[lowestCacheHitRateId]
		addCacheRow(id, textToAddToRow)
		pushToDisk(lowestCacheHitRateId, rowToPushToMemory.text)
		deleteRowFromCache(lowestCacheHitRateId)
	}
}

func getLowestHitRowId() uint64 {
	var lowestHitRateId uint64 = MAX_UINT64
	var lowestHitRateRow CacheRow
	for id, cacheRow := range cache {
		if cacheRow.selectCount < lowestHitRateRow.selectCount {
			lowestHitRateId = id
			lowestHitRateRow = cacheRow
		}
	}
	return lowestHitRateId
}

func deleteRowFromCache(id uint64) {
	_, canDel := cache[id];
    if canDel {
		delete(cache, id)
    }
}

func addCacheRow(id uint64, text string) {
	cache[id] = CacheRow{text, 0, false}
}


// Find in Cache or get from memory
func findRow(id uint64) {
	cacheRow, textFound := cache[id]
	if (textFound) {
		cacheRow.selectCount++
		fmt.Printf("%d: %s\n", id, cacheRow.text)
	} else {
		memoryRow, found := getRowFromDisk(id)
		if found {
			if len(cache) == MAX_CACHE_SIZE {
				lowestCacheHitRateId := getLowestHitRowId()
				deleteRowFromCache(lowestCacheHitRateId)
			}
			addMemoryRowToCacheTable(memoryRow)
			fmt.Printf("%d: %s\n", memoryRow.id, memoryRow.text)
		} else {
			print("no texts found\n")			
		}
	}
}

func addMemoryRowToCacheTable(memoryRow TextDataRow) {
	cache[memoryRow.id] = CacheRow{memoryRow.text,1,true}	
}

/*
	on insert
push to cache
if cache is full push: push least selected to memory

	on exit
push all cache not already in memory to memory

	on select
check cache
if not in cache check memory
if not in memory return nothing

	load cache
get up to max cache size from disk
set all as in memory
*/
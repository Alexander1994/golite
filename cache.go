package main

import "math"

const maxRowBytes = 1024
const maxCacheSize = 20

type cacheRow struct {
	text        string
	selectCount uint32
	inMem       bool
}

var cache = make(map[uint32]cacheRow)

func getLowestHitRowID() uint32 {
	var lowestHitRateID uint32 = math.MaxUint32
	var lowestHitRateRow cacheRow
	for id, cacheRow := range cache {
		if cacheRow.selectCount <= lowestHitRateRow.selectCount {
			lowestHitRateID = id
			lowestHitRateRow = cacheRow
		}
	}
	return lowestHitRateID
}

func deleteRowFromCache(id uint32) {
	_, canDel := cache[id]
	if canDel {
		delete(cache, id)
	}
}

func addCacheRow(id uint32, text string) {
	cache[id] = cacheRow{text, 0, false}
}

func addMemRowToCache(id uint32, text string) {
	addCacheRow(id, text)
	cacheRow := cache[id]
	cacheRow.selectCount++
}

func resetCache() {
	cache = make(map[uint32]cacheRow)
}

func closeCache() {
	for id, cacheRow := range cache {
		if !cacheRow.inMem {
			pushToDisk(id, cacheRow.text)
		}
	}
}

package main

import "math"

const maxRowBytes = 1024
const maxCacheSize = 20

type cacheRow struct {
	text        string
	selectCount uint64
	inMem       bool
}

var cache = make(map[uint64]cacheRow)

func getLowestHitRowID() uint64 {
	var lowestHitRateID uint64 = math.MaxUint64
	var lowestHitRateRow cacheRow
	for id, cacheRow := range cache {
		if cacheRow.selectCount < lowestHitRateRow.selectCount {
			lowestHitRateID = id
			lowestHitRateRow = cacheRow
		}
	}
	return lowestHitRateID
}

func deleteRowFromCache(id uint64) {
	_, canDel := cache[id]
	if canDel {
		delete(cache, id)
	}
}

func addCacheRow(id uint64, text string) {
	cache[id] = cacheRow{text, 0, false}
}

func addMemoryRowToCacheTable(memoryRow TextDataRow) {
	cache[memoryRow.id] = cacheRow{memoryRow.text, 1, true}
}

func resetCache() {
	cache = make(map[uint64]cacheRow)
}

package database

import "math"

const maxRowBytes = 1024
const maxCacheSize = 20

type cacheRow struct {
	text        string
	selectCount uint32
	inMem       bool
}

type cacheTable map[uint32]cacheRow

func createCache() cacheTable {
	return make(cacheTable)
}

func (cache cacheTable) getLowestHitRowID() uint32 {
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

func (cache cacheTable) deleteRow(id uint32) {
	_, canDel := cache[id]
	if canDel {
		delete(cache, id)
	}
}

func (cache cacheTable) addRow(id uint32, text string) {
	cache[id] = cacheRow{text, 0, false}
}

func (cache cacheTable) addMemRow(id uint32, text string) {
	cache.addRow(id, text)
	cacheRow := cache[id]
	cacheRow.selectCount++
}

func (cache cacheTable) reset() {
	cache = make(map[uint32]cacheRow)
}

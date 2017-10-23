package main

import (
	"math"
	"sort"
)

type pageRow struct {
	offset uint32
	length uint16
}

var pgTable []pageRow

func insertRowIntoPageTable(offset uint32, length uint16) {
	currRange := pageRow{offset, length}
	pgTable = append(pgTable, currRange)
}

func getSmallestHoleToFit(length uint16, nextMetaTableOffset uint32) (uint32, bool) { // offset relative to meta table end, offset found
	orderPgTable()
	pgTableLen := len(pgTable)
	smallestHoleOffset := uint32(math.MaxUint32)
	currHoleSize := uint32(0)
	foundHole := false

	if len(pgTable) == 0 {
		return 0, true
	}
	for i := 0; i < pgTableLen; i++ {
		ithEndOffset := getIthEndOffset(i)
		if nextMetaTableOffset == 0 && i == len(pgTable)-1 && !foundHole { // if no next meta table escape, no holes & no more pgs
			return getIthEndOffset(len(pgTable) - 1), true // return end offset
		}
		currHoleSize = getHoleSize(i, nextMetaTableOffset)
		if uint16(currHoleSize) == length { // if perfect hole use
			return ithEndOffset, true
		}
		if currHoleSize > uint32(length) && currHoleSize < smallestHoleOffset { // current hole is smallest and fits text
			smallestHoleOffset = ithEndOffset
			foundHole = true
		}
	}
	if foundHole {
		return smallestHoleOffset, true
	}
	return 0, false
}

func getHoleSize(i int, nextMetaTableOffset uint32) uint32 {
	ithEndOffset := getIthEndOffset(i)
	if i == len(pgTable)-1 { // if last index use space to next meta table
		return nextMetaTableOffset - ithEndOffset
	}
	return pgTable[i+1].offset - ithEndOffset
}

func resetPgTable() {
	pgTable = nil
}

func orderPgTable() {
	sort.Slice(pgTable, func(i, j int) bool {
		return pgTable[i].offset < pgTable[j].offset
	})
}

func getIthEndOffset(i int) uint32 {
	return pgTable[i].offset + uint32(pgTable[i].length)
}

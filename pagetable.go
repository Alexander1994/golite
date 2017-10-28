package main

import (
	"math"
	"sort"
)

type pageRow struct {
	offset uint32
	length uint16
}
type pageTable []pageRow

var pgTable pageTable

func (pgTable pageTable) insertRow(offset uint32, length uint16) {
	currRange := pageRow{offset, length}
	pgTable = append(pgTable, currRange)
}

func (pgTable pageTable) getSmallestHoleToFit(length uint16, nextMetaTableOffset uint32) (uint32, bool) { // offset relative to meta table end, offset found
	pgTable.orderByOffset()
	pgTableLen := len(pgTable)
	smallestHoleOffset := uint32(math.MaxUint32)
	currHoleSize := uint32(0)
	foundHole := false

	if len(pgTable) == 0 {
		return 0, true
	}
	for i := 0; i < pgTableLen; i++ {
		ithEndOffset := pgTable.getIthEndOffset(i)
		if nextMetaTableOffset == 0 && i == len(pgTable)-1 && !foundHole { // if no next meta table escape, no holes & no more pgs
			return pgTable.getIthEndOffset(len(pgTable) - 1), true // return end offset
		}
		currHoleSize = pgTable.getHoleSize(i, nextMetaTableOffset)
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

func (pgTable pageTable) getHoleSize(i int, nextMetaTableOffset uint32) uint32 {
	ithEndOffset := pgTable.getIthEndOffset(i)
	if i == len(pgTable)-1 { // if last index use space to next meta table
		return nextMetaTableOffset - ithEndOffset
	}
	return pgTable[i+1].offset - ithEndOffset
}

func (pgTable pageTable) reset() {
	pgTable = nil
}

func (pgTable pageTable) orderByOffset() {
	sort.Slice(pgTable, func(i, j int) bool {
		return pgTable[i].offset < pgTable[j].offset
	})
}

func (pgTable pageTable) getIthEndOffset(i int) uint32 {
	return pgTable[i].offset + uint32(pgTable[i].length)
}

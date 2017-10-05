package main

import (
	"math"
	"sort"
)

type PageRow struct {
	offset int64
	length int32
}

var pageTable []PageRow // length, offset

func loadPageTable() {
	if size == 0 {
		return
	}
	textLengthByteArr := make([]byte, textLengthByteLength)

	for {
		offset := currentOffSet()
		spaceLength, found := seekOverSpaceToId()
		if !found {
			break
		}
		if spaceLength != 0 {
			pageTable = append(pageTable, PageRow{offset, spaceLength})
		}
		if !seekOverRow(textLengthByteArr) {
			break
		}
	}
}

// find smallest offset that fits length
func getInsertOffset(toInsertLength int32) (int64, bool) {
	currSmallestLength := int32(math.MaxInt32) // max int, for comparison
	var offset int64
	found := false
	for _, row := range pageTable {
		if row.length >= toInsertLength && row.length < currSmallestLength {
			currSmallestLength = row.length
			offset = row.offset
			found = true
		}
	}
	return offset, found
}

func addPageToPageTable(offset int64, textLength uint16) {
	pageRow := PageRow{offset, rowByteLength(textLength)}
	pageTable = append(pageTable, pageRow)
	updatePageTable()
}

func removePageFromPageTable(i int) {
	pageTable = append(pageTable[:i], pageTable[i+1:]...)
}

func resetPageTable() {
	pageTable = nil
}

// When offset in page table & value has been updated and page table needs to be updated.
func removeLengthFromOffset(offset int64, rowLength int32) {
	for _, row := range pageTable {
		if row.offset == offset {
			if rowLength <= row.length {
				row.length -= rowLength
			}
		}
	}
	updatePageTable()
}

// For sorting page row by offset
type ByOffset []PageRow

func (a ByOffset) Len() int           { return len(a) }
func (a ByOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOffset) Less(i, j int) bool { return a[i].offset < a[j].offset }

// merge overlapping offset + lengths in page table
func updatePageTable() {
	sort.Sort(ByOffset(pageTable))
	for i := 1; i < len(pageTable); i++ {
		if pageTable[i-1].offset+int64(pageTable[i-1].length) == pageTable[i].offset {
			pageTable[i].offset = pageTable[i-1].offset
			pageTable[i].length += pageTable[i-1].length
			removePageFromPageTable(i - 1)
		}
	}
}

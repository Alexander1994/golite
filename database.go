package main

import (
	"fmt"
	"log"
	"os"
)

const dirname = ".data"

type TextDataRow struct {
	id         uint32
	textLength uint16
	text       string
}

var (
	file    *os.File
	fileErr error
)

var size int64
var fileName = dirname + "/db.dat"

// DB controls
func openDisk(testMode bool) {
	if testMode {
		fileName = dirname + "/testdb.dat"
	}
	os.Mkdir(dirname, 0755)

	file, fileErr = os.OpenFile(fileName,
		os.O_RDWR|os.O_CREATE,
		0600)
	fatal(fileErr)
	fileStat, _ := file.Stat()
	size = fileStat.Size()
	if size < int64(osPageSize) {
		emptyMetaTableRow := make([]byte, osPageSize)
		file.Write(emptyMetaTableRow)
	}
}

func closeDisk() {
	if file != nil {
		file.Close()
	} else {
		println("open db before attempting to close it")
	}
}

// DB commands
func pushToDisk(id uint32, text string) bool {
	nextMetaTableOffset := uint32(0)
	textLength := uint16(len(text))

	resetCursorToStart()

	for true {
		loadMetaTable(nextMetaTableOffset)
		insertLocation := metaTableMaxRowCount
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := getID(i)
			if id == ithID {
				return false
			}
			if ithID == 0 && insertLocation == metaTableMaxRowCount { // ID = empty & insert location is not found
				insertLocation = i
			}
			if ithID != 0 {
				insertRowIntoPageTable(getTextOffset(i), getLength(i))
			}
		}
		if insertLocation != metaTableMaxRowCount { // if hole was found
			nextMetaTableOffset = getMetaTableOffset()
			offset, found := getSmallestHoleToFit(textLength, nextMetaTableOffset)
			resetPgTable()
			if found {
				setMetaTableRow(insertLocation, id, textLength, offset)
				setTextRow(offset, text)
				return true
			}
		} else { // if no holes found add next table and seek to
			addAndGoToNextMetaTable()
			nextMetaTableOffset = 0 // since we are at the next metatable offset=0
			resetPgTable()
		}
	}
	return false // should never happen, throw err?
}

func getRowFromDisk(id uint32) (string, bool) { // text, found
	resetCursorToStart()
	nextMetaTableOffset := uint32(0)
	for true {
		loadMetaTable(nextMetaTableOffset)
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := getID(i)
			if id == ithID {
				return getText(getTextOffset(i), getLength(i)), true
			}
		}

		nextMetaTableOffset = getMetaTableOffset()
		if nextMetaTableOffset == 0 {
			return "", false
		}
	}
	return "", false
}

func deleteRowFromDisk(id uint32) bool {
	resetCursorToStart()
	nextMetaTableOffset := uint32(0)
	for true {
		loadMetaTable(nextMetaTableOffset)
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := getID(i)
			if id == ithID {
				deleteIthRow(i)
				return true
			}
		}
		nextMetaTableOffset = getMetaTableOffset()
		if nextMetaTableOffset == 0 {
			return false
		}
	}
	return false
}

// Cursor info & controls
func currentOffSet() int64 {
	offset, e := file.Seek(0, 1)
	fatal(e)
	return offset
}

func resetCursorToStart() {
	file.Seek(0, 0)
}

func printCursorOffset() {
	fmt.Printf("cursor offset:%d\n", currentOffSet())
}

// Logging tool
func fatal(err error) {
	if err != nil {
		print("\n")
		log.Fatal(err)
		print("\n")
	}
}

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

type disk os.File

var (
	file    *disk
	fileErr error
)

var size int64
var fileName = dirname + "/db.dat"

// DB controls
func (file *disk) open(testMode bool) {
	if testMode {
		fileName = dirname + "/testdb.dat"
	}
	os.Mkdir(dirname, 0755)

	f, err := os.OpenFile(fileName,
		os.O_RDWR|os.O_CREATE,
		0600)
	fatal(err)
	file = (*disk)(f)
	if size < int64(osPageSize) {
		emptyMetaTableRow := make([]byte, osPageSize)
		file.write(emptyMetaTableRow)
	}
}

func (file *disk) close() {
	if file != nil {
		(*os.File)(file).Close()
	} else {
		println("open db before attempting to close it")
	}
}

// DB commands
func (file *disk) pushToDisk(id uint32, text string) bool {
	nextMetaTableOffset := uint32(0)
	textLength := uint16(len(text))

	file.resetCursorToStart()

	for true {
		metaTable := file.loadMetaTable(nextMetaTableOffset)
		insertLocation := metaTableMaxRowCount
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := metaTable.getID(i)
			if id == ithID {
				return false
			}
			if ithID == 0 && insertLocation == metaTableMaxRowCount { // ID = empty & insert location is not found
				insertLocation = i
			}
			if ithID != 0 {
				pgTable.insertRow(metaTable.getTextOffset(i), metaTable.getLength(i))
			}
		}
		if insertLocation != metaTableMaxRowCount { // if hole was found
			nextMetaTableOffset = metaTable.getMetaTableOffset()
			offset, found := pgTable.getSmallestHoleToFit(textLength, nextMetaTableOffset)
			pgTable.reset()
			if found {
				file.setMetaTableRow(insertLocation, id, textLength, offset)
				file.setTextRow(offset, text)
				return true
			}
		} else { // if no holes found add next table and seek to
			file.addAndGoToNextMetaTable()
			nextMetaTableOffset = 0 // since we are at the next metatable offset=0
			pgTable.reset()
		}
	}
	return false // should never happen, throw err?
}

func (file *disk) getRowFromDisk(id uint32) (string, bool) { // text, found
	file.resetCursorToStart()
	nextMetaTableOffset := uint32(0)
	for true {
		metaTable := file.loadMetaTable(nextMetaTableOffset)
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := metaTable.getID(i)
			if id == ithID {
				return file.getText(metaTable.getTextOffset(i), metaTable.getLength(i)), true
			}
		}

		nextMetaTableOffset = metaTable.getMetaTableOffset()
		if nextMetaTableOffset == 0 {
			return "", false
		}
	}
	return "", false
}

func (file *disk) deleteRowFromDisk(id uint32) bool {
	file.resetCursorToStart()
	nextMetaTableOffset := uint32(0)
	for true {
		metaTable := file.loadMetaTable(nextMetaTableOffset)
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := metaTable.getID(i)
			if id == ithID {
				file.deleteIthRow(i)
				return true
			}
		}
		nextMetaTableOffset = metaTable.getMetaTableOffset()
		if nextMetaTableOffset == 0 {
			return false
		}
	}
	return false
}

// Cursor info & controls
func (file *disk) seek(offset int64) {
	_, err := (*os.File)(file).Seek(offset, 1)
	fatal(err)
}

func (file *disk) write(b []byte) {
	_, err := (*os.File)(file).Write(b)
	fatal(err)
}

func (file *disk) read(b []byte) {
	_, err := (*os.File)(file).Read(b)
	fatal(err)
}

func (file *disk) size() int64 {
	fileStat, err := ((*os.File)(file)).Stat()
	fatal(err)
	return fileStat.Size()
}

func (file *disk) currentOffSet() int64 {
	offset, e := ((*os.File)(file)).Seek(0, 1)
	fatal(e)
	return offset
}

func (file *disk) resetCursorToStart() {
	(*os.File)(file).Seek(0, 0)
}

func (file *disk) printCursorOffset() {
	fmt.Printf("cursor offset:%d\n", file.currentOffSet())
}

// Logging tool
func fatal(err error) {
	if err != nil {
		print("\n")
		log.Fatal(err)
		print("\n")
	}
}

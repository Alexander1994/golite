package database

import (
	"os"
)

type TextDataRow struct {
	id         uint32
	textLength uint16
	text       string
}

type disk os.File

// DB controls
func createDisk(fileName string) *disk {
	f, err := os.OpenFile(dirname+"/"+fileName+".dat",
		os.O_RDWR|os.O_CREATE,
		0600)
	fatal(err)
	file := (*disk)(f)
	if file.size() < int64(osPageSize) {
		emptyMetaTableRow := make([]byte, osPageSize)
		file.write(emptyMetaTableRow)
	}
	return file
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
	var pgTable pageTable

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
				pgTable = pgTable.insertRow(metaTable.getTextOffset(i), metaTable.getLength(i))
			}
		}
		if insertLocation != metaTableMaxRowCount { // if hole was found
			nextMetaTableOffset = metaTable.getMetaTableOffset()
			offset, found := pgTable.getSmallestHoleToFit(textLength, nextMetaTableOffset)
			if found {
				file.setMetaTableRow(insertLocation, id, textLength, offset)
				file.setTextRow(offset, text)
				return true
			}
		} else if metaTable.getMetaTableOffset() == 0 { // if no holes found add next table and seek to
			nextMetaTableOffset := metaTable.getTextOffset(metaTableMaxRowCount-1) + uint32(metaTable.getLength(metaTableMaxRowCount-1))
			file.addAndGoToNextMetaTable(nextMetaTableOffset)
			nextMetaTableOffset = 0 // since we are at the next metatable offset=0
		} else { // table full and is next offset
			nextMetaTableOffset = metaTable.getMetaTableOffset()
		}
	}
	return false // should never happen, throw err?
}

func (file *disk) getAllRowsFromDisk() []TextDataRow {
	file.resetCursorToStart()
	var allRows []TextDataRow
	nextMetaTableOffset := uint32(0)
	for true {
		metaTable := file.loadMetaTable(nextMetaTableOffset)
		for i := uint32(0); i < metaTableMaxRowCount; i++ {
			ithID := metaTable.getID(i)
			ithTextLen := metaTable.getLength(i)
			ithTextOffset := metaTable.getTextOffset(i)
			ithText := file.getText(ithTextOffset, ithTextLen)
			allRows = append(allRows, TextDataRow{ithID, ithTextLen, ithText})
			file.seek(-int64(ithTextOffset + uint32(ithTextLen)))
		}
		nextMetaTableOffset = metaTable.getMetaTableOffset()
		if nextMetaTableOffset == 0 {
			return allRows
		}
	}
	return allRows
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

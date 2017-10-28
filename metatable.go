package main

import (
	"encoding/binary"
	"os"
)

const (
	idByteLen              = 4 // uint32
	textLenByteLen         = 2 // uint16
	offsetByteLen          = 4 // uint32
	metaRowByteLen         = idByteLen + textLenByteLen + offsetByteLen
	metaTableOffsetByteLen = 4
)

// MetaTable ...
type MetaTable []byte

var (
	osPageSize                              = int32(os.Getpagesize())
	metaTable                     MetaTable = make([]byte, osPageSize)
	metaTableMaxRowCount                    = uint32((osPageSize - metaTableOffsetByteLen) / metaRowByteLen)
	metaTableOffsetAndXtraByteLen           = uint32(osPageSize) - (metaTableMaxRowCount * metaRowByteLen)
)

// Meta Table getters
func (metaTable MetaTable) getID(ithRow uint32) uint32 {
	offset := ithRow * metaRowByteLen
	return binary.LittleEndian.Uint32(metaTable[offset : idByteLen+offset])
}

func (metaTable MetaTable) getLength(ithRow uint32) uint16 {
	offset := ithRow*metaRowByteLen + idByteLen
	return binary.LittleEndian.Uint16(metaTable[offset : textLenByteLen+offset])
}

func (metaTable MetaTable) getTextOffset(ithRow uint32) uint32 {
	offset := ithRow*metaRowByteLen + idByteLen + textLenByteLen
	return binary.LittleEndian.Uint32(metaTable[offset : offsetByteLen+offset])
}

func (metaTable MetaTable) getMetaTableOffset() uint32 {
	return binary.LittleEndian.Uint32(metaTable[osPageSize-metaTableOffsetByteLen:])
}

func (file *disk) getText(offset uint32, length uint16) string {
	file.seek(int64(offset))
	textByteArr := make([]byte, length)
	file.read(textByteArr)
	return string(textByteArr)
}

// Text data setter
func (file *disk) setTextRow(offset uint32, text string) {
	file.seek(int64(offset))
	((*os.File)(file)).WriteString(text)
}

// Meta Table setter
func (file *disk) setMetaTableRow(ithRow uint32, id uint32, textLength uint16, offset uint32) {
	idByteArr := make([]byte, idByteLen)
	textLengthByteArr := make([]byte, textLenByteLen)
	offsetByteArr := make([]byte, offsetByteLen)
	var rowByteArr []byte

	binary.LittleEndian.PutUint32(idByteArr, id)
	binary.LittleEndian.PutUint16(textLengthByteArr, textLength)
	binary.LittleEndian.PutUint32(offsetByteArr, offset)

	rowByteArr = append(rowByteArr, idByteArr...)
	rowByteArr = append(rowByteArr, textLengthByteArr...)
	rowByteArr = append(rowByteArr, offsetByteArr...)

	file.writeInIthRowAndSeekToTableEnd(ithRow, rowByteArr)
}

// Delete ith meta table row
func (file *disk) deleteIthRow(i uint32) { // bytes seeked to reverse to end
	emptyRow := make([]byte, metaRowByteLen)
	seekAmount := metaTableOffsetAndXtraByteLen + (metaTableMaxRowCount-i)*metaRowByteLen
	file.seek(-int64(seekAmount)) // seek back to index
	file.write(emptyRow)
}

// load meta table into memory from disk
func (file *disk) loadMetaTable(offset uint32) MetaTable {
	metaTable = make(MetaTable, osPageSize)
	if offset > 0 {
		file.seek(int64(offset))
	}
	file.read(metaTable)
	return metaTable
}

// add a new meta table and set cursor offset to next meta table start
func (file *disk) addAndGoToNextMetaTable() {
	lastPgIndex := uint32(len(pgTable) - 1)
	nextMetaTableOffset := metaTable.getTextOffset(lastPgIndex) + uint32(metaTable.getLength(lastPgIndex))
	file.seek(-int64(metaTableOffsetByteLen)) // seek to current metatable offset location and set next metatable offset
	nextMetaTableOffsetByteArr := make([]byte, metaTableOffsetByteLen)
	binary.LittleEndian.PutUint32(nextMetaTableOffsetByteArr, nextMetaTableOffset)
	file.write(nextMetaTableOffsetByteArr)

	file.seek(int64(nextMetaTableOffset)) // seek to next metatable

	emptyTable := make([]byte, osPageSize) // create metatable on file
	file.write(emptyTable)

	file.seek(-int64(osPageSize)) // Seek and adjust to start before metatable
}

func (file *disk) writeInIthRowAndSeekToTableEnd(i uint32, row []byte) { // bytes seeked to reverse to end
	seekAmount := metaTableOffsetAndXtraByteLen + (metaTableMaxRowCount-i)*metaRowByteLen
	file.seek(-int64(seekAmount)) // seek back to index
	file.write(row)
	file.seek(int64(seekAmount - uint32(len(row))))
}

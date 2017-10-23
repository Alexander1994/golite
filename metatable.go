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

var (
	osPageSize                    = int32(os.Getpagesize())
	metaTable                     = make([]byte, osPageSize)
	metaTableMaxRowCount          = uint32((osPageSize - metaTableOffsetByteLen) / metaRowByteLen)
	metaTableOffsetAndXtraByteLen = uint32(osPageSize) - (metaTableMaxRowCount * metaRowByteLen)
)

// Meta Table getters
func getID(ithRow uint32) uint32 {
	offset := ithRow * metaRowByteLen
	return binary.LittleEndian.Uint32(metaTable[offset : idByteLen+offset])
}

func getLength(ithRow uint32) uint16 {
	offset := ithRow*metaRowByteLen + idByteLen
	return binary.LittleEndian.Uint16(metaTable[offset : textLenByteLen+offset])
}

func getTextOffset(ithRow uint32) uint32 {
	offset := ithRow*metaRowByteLen + idByteLen + textLenByteLen
	return binary.LittleEndian.Uint32(metaTable[offset : offsetByteLen+offset])
}

func getMetaTableOffset() uint32 {
	return binary.LittleEndian.Uint32(metaTable[osPageSize-metaTableOffsetByteLen:])
}

func getText(offset uint32, length uint16) string {
	file.Seek(int64(offset), 1)
	textByteArr := make([]byte, length)
	file.Read(textByteArr)
	return string(textByteArr)
}

// Text data setter
func setTextRow(offset uint32, text string) {
	_, err := file.Seek(int64(offset), 1)
	fatal(err)
	file.WriteString(text)
}

// Meta Table setter
func setMetaTableRow(ithRow uint32, id uint32, textLength uint16, offset uint32) {
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

	writeInIthRowAndSeekToTableEnd(ithRow, rowByteArr)
}

// Delete ith meta table row
func deleteIthRow(i uint32) { // bytes seeked to reverse to end
	emptyRow := make([]byte, metaRowByteLen)
	seekAmount := metaTableOffsetAndXtraByteLen + (metaTableMaxRowCount-i)*metaRowByteLen
	file.Seek(-int64(seekAmount), 1) // seek back to index
	_, err := file.Write(emptyRow)
	fatal(err)
}

// load meta table into memory from disk
func loadMetaTable(offset uint32) {
	metaTable = make([]byte, osPageSize)
	if offset > 0 {
		_, err := file.Seek(int64(offset), 1)
		fatal(err)
	}
	_, err := file.Read(metaTable)
	fatal(err)
}

// add a new meta table and set cursor offset to next meta table start
func addAndGoToNextMetaTable() {
	lastPgIndex := uint32(len(pgTable) - 1)
	nextMetaTableOffset := getTextOffset(lastPgIndex) + uint32(getLength(lastPgIndex))
	file.Seek(-int64(metaTableOffsetByteLen), 1) // seek to current metatable offset location and set next metatable offset
	nextMetaTableOffsetByteArr := make([]byte, metaTableOffsetByteLen)
	binary.LittleEndian.PutUint32(nextMetaTableOffsetByteArr, nextMetaTableOffset)
	file.Write(nextMetaTableOffsetByteArr)

	file.Seek(int64(nextMetaTableOffset), 1) // seek to next metatable

	emptyTable := make([]byte, osPageSize) // create metatable on file
	file.Write(emptyTable)

	file.Seek(-int64(osPageSize), 1) // Seek and adjust to start before metatable
}

func writeInIthRowAndSeekToTableEnd(i uint32, row []byte) { // bytes seeked to reverse to end
	seekAmount := metaTableOffsetAndXtraByteLen + (metaTableMaxRowCount-i)*metaRowByteLen
	file.Seek(-int64(seekAmount), 1) // seek back to index
	_, err := file.Write(row)
	fatal(err)
	file.Seek(int64(seekAmount-uint32(len(row))), 1)
}

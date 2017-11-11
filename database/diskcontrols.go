package database

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

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

// Delete ith meta table row from disk
func (file *disk) deleteIthRow(i uint32) { // bytes seeked to reverse to end
	emptyRow := make([]byte, metaRowByteLen)
	seekAmount := metaTableOffsetAndXtraByteLen + (metaTableMaxRowCount-i)*metaRowByteLen
	file.seek(-int64(seekAmount)) // seek back to index
	file.write(emptyRow)
}

// load meta table into memory from disk
func (file *disk) loadMetaTable(offset uint32) MetaTable {
	metaTable := make(MetaTable, osPageSize)
	if offset > 0 {
		file.seek(int64(offset))
	}
	file.read(metaTable)
	return metaTable
}

// add a new meta table and set cursor offset to next meta table start
func (file *disk) addAndGoToNextMetaTable(nextMetaTableOffset uint32) {
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

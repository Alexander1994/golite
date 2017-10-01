package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

const ID_BYTE_LENGTH = 8
const TEXT_LENGTH_BYTE_LENGTH = 2

const DIRNAME = ".data"
const FILENAME = DIRNAME + "/db.dat"

/*
 *   63 bits | 16 bit length | var bit length, max length 65536 *note zero length not option, 1 bit for identification
 *      ID   |  textLength   | text
 */
type TextDataRow struct {
	id         uint64
	textLength uint16
	text       string
}

var (
	file *os.File
	err  error
)

var size int64

// DB controls
func openDB() {
	os.Mkdir(DIRNAME, 0755)

	file, err = os.OpenFile(FILENAME,
		os.O_RDWR|os.O_CREATE,
		0600)
	panic(err)
	fileStat, _ := file.Stat()
	size = fileStat.Size()
}

func closeDB() {
	if file != nil {
		for id, cacheRow := range cache {
			if !cacheRow.inMem {
				pushToDisk(id, cacheRow.text)
			}
		}
	} else {
		println("open db first.")
	}
	file.Close()
}

// DB commands
func pushToDisk(id uint64, text string) {
	textLength := uint16(len(text))
	rowToPush := TextDataRow{idToDBID(id), textLength, text}
	byteArr := textDataRowToBytes(rowToPush)
	resetCursorToStart()
	if size == 0 {
		file.Write(byteArr)
	} else {
		insertOffset, found := getInsertOffset(rowByteLength(textLength))

		if found {
			removeLengthFromOffset(insertOffset, rowByteLength(textLength))
			_, err := file.WriteAt(byteArr, insertOffset)
			panic(err)
		} else {
			_, err := file.WriteAt(byteArr, size)
			panic(err)
		}
	}
	fileStat, _ := file.Stat()
	size = fileStat.Size()
}

func getRowFromDisk(id uint64) (TextDataRow, bool) {
	textLengthByteArr := make([]byte, TEXT_LENGTH_BYTE_LENGTH)
	idByteArr := make([]byte, ID_BYTE_LENGTH)
	resetCursorToStart()
	for {
		_, found := seekOverSpaceToId()
		if !found {
			break
		}
		rowId, found := readId(idByteArr)
		if !found {
			break
		}
		textLength, found := readTextLength(textLengthByteArr)
		if !found {
			break
		}
		if rowId == id {
			text := readText(textLength)
			return TextDataRow{rowId, textLength, text}, true
		}
		if !seekOverText(textLength) {
			break
		}
	}

	return TextDataRow{}, false
}

func deleteRowFromDisk(id uint64) bool {
	textLengthByteArr := make([]byte, TEXT_LENGTH_BYTE_LENGTH)
	idByteArr := make([]byte, ID_BYTE_LENGTH)
	resetCursorToStart()
	for {
		_, found := seekOverSpaceToId()
		if !found {
			break
		}
		offset := currentOffSet()
		rowId, found := readId(idByteArr)
		if !found {
			break
		}

		textLength, found := readTextLength(textLengthByteArr)
		if !found {
			break
		}
		if rowId == id {
			addPageToPageTable(offset, textLength)
			deleteRowInDisk(textLength)
			fileStat, _ := file.Stat()
			size = fileStat.Size()
			return true
		}
		if !seekOverText(textLength) {
			break
		}
	}

	return false
}

// Read functions
func readId(idByteArr []byte) (uint64, bool) {
	bytesRead, e := file.Read(idByteArr)
	panic(e)
	if bytesRead == 0 {
		return 0, false
	}
	test := binary.BigEndian.Uint64(idByteArr)
	return dbIDToID(test), true
}

func readTextLength(textLengthByteArr []byte) (uint16, bool) {
	bytesRead, e := file.Read(textLengthByteArr)
	panic(e)
	if bytesRead == 0 {
		return 0, false
	}
	return binary.BigEndian.Uint16(textLengthByteArr), true
}

func readText(textLength uint16) string {
	textByteArr := make([]byte, int(textLength))
	_, e := file.Read(textByteArr)
	panic(e)
	return string(textByteArr)
}

// Seek functions
func seekOverText(textLength uint16) bool {
	seekLength, e := file.Seek(int64(textLength), 1)
	panic(e)
	return seekLength < size

}

func seekOverID() bool {
	seekLength, e := file.Seek(ID_BYTE_LENGTH, 1)
	panic(e)
	return seekLength < size
}

func seekOverRow(textLengthByteArr []byte) bool {
	if !seekOverID() {
		return false
	}
	textLength, found := readTextLength(textLengthByteArr)
	if !found {
		return false
	}
	if !seekOverText(textLength) {
		return false
	}
	return true
}

func seekOverSpaceToId() (int32, bool) {
	var emptyByteCount int32
	byteArr := []byte{0}
	emptyArr := []byte{0}

	for {
		_, err := file.Read(byteArr)
		if err == io.EOF {
			return emptyByteCount, false
		}

		if !bytes.Equal(byteArr, emptyArr) {
			_, e := file.Seek(-1, 1)
			panic(e)
			return emptyByteCount, true
		}
		emptyByteCount++

	}
}

// conversions
func textDataRowToBytes(row TextDataRow) []byte {
	idByteArr := make([]byte, ID_BYTE_LENGTH)
	textLengthByteArr := make([]byte, TEXT_LENGTH_BYTE_LENGTH)
	var rowByteArr []byte

	textByteArr := []byte(row.text)
	binary.BigEndian.PutUint64(idByteArr, row.id)
	binary.BigEndian.PutUint16(textLengthByteArr, row.textLength)

	rowByteArr = append(rowByteArr, idByteArr...)
	rowByteArr = append(rowByteArr, textLengthByteArr...)
	rowByteArr = append(rowByteArr, textByteArr...)

	return rowByteArr
}

func bytesToTextDataRow(b []byte) TextDataRow {
	id := binary.BigEndian.Uint64(b[:8])
	len := binary.BigEndian.Uint16(b[8:10])
	text := string(b[10 : len+10])
	return TextDataRow{id, len, text}
}

// helper functions
func deleteRowInDisk(textLength uint16) {
	_, e := file.Seek(-(TEXT_LENGTH_BYTE_LENGTH + ID_BYTE_LENGTH), 1)
	panic(e)
	nullRow := make([]byte, rowByteLength(textLength))
	file.Write(nullRow)
}

func currentOffSet() int64 {
	offset, e := file.Seek(0, 1)
	panic(e)
	return offset
}

func resetCursorToStart() {
	file.Seek(0, 0)
}

// for debug purposes
func panic(err error) {
	if err != nil {
		print("\n")
		log.Fatal(err)
		print("\n")
	}
}

func rowByteLength(textLength uint16) int32 {
	return int32(TEXT_LENGTH_BYTE_LENGTH + ID_BYTE_LENGTH + textLength)
}

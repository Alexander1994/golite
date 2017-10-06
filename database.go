package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

const idByteLength = 8
const textLengthByteLength = 2

const dirname = ".data"

/*
 *  1 bite        | 63 bits | 16 bit length | var bit length, max length 65536 *note zero length not option, 1 bit for identification
 * row identifier | ID      | textLength    | text
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
var fileName = dirname + "/db.dat"

// DB controls
func openDisk(testMode bool) {
	if testMode {
		fileName = dirname + "/testdb.dat"
	}
	os.Mkdir(dirname, 0755)

	file, err = os.OpenFile(fileName,
		os.O_RDWR|os.O_CREATE,
		0600)
	panic(err)
	fileStat, _ := file.Stat()
	size = fileStat.Size()
}

func closeDisk() {
	if file != nil {
		file.Close()
	} else {
		println("open db before attempting to close it")
	}
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
	textLengthByteArr := make([]byte, textLengthByteLength)
	idByteArr := make([]byte, idByteLength)
	resetCursorToStart()
	for {
		_, found := seekOverSpaceToID()
		if !found {
			break
		}
		rowID, found := readID(idByteArr)
		if !found {
			break
		}
		textLength, found := readTextLength(textLengthByteArr)
		if !found {
			break
		}
		if rowID == id {
			text := readText(textLength)
			return TextDataRow{rowID, textLength, text}, true
		}
		if !seekOverText(textLength) {
			break
		}
	}

	return TextDataRow{}, false
}

func deleteRowFromDisk(id uint64) bool {
	textLengthByteArr := make([]byte, textLengthByteLength)
	idByteArr := make([]byte, idByteLength)
	resetCursorToStart()
	for {
		_, found := seekOverSpaceToID()
		if !found {
			break
		}
		offset := currentOffSet()
		rowID, found := readID(idByteArr)
		if !found {
			break
		}

		textLength, found := readTextLength(textLengthByteArr)
		if !found {
			break
		}
		if rowID == id {
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
func readID(idByteArr []byte) (uint64, bool) {
	bytesRead, e := file.Read(idByteArr)
	panic(e)
	if bytesRead == 0 {
		return 0, false
	}
	dbID := binary.BigEndian.Uint64(idByteArr)
	return dbIDToID(dbID), true
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
	seekLength, e := file.Seek(idByteLength, 1)
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

func seekOverSpaceToID() (int32, bool) {
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
	idByteArr := make([]byte, idByteLength)
	textLengthByteArr := make([]byte, textLengthByteLength)
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
	_, e := file.Seek(-(textLengthByteLength + idByteLength), 1)
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

func rowByteLength(textLength uint16) int32 {
	return int32(textLengthByteLength + idByteLength + textLength)
}

// for debug purposes
func panic(err error) {
	if err != nil {
		print("\n")
		log.Fatal(err)
		print("\n")
	}
}

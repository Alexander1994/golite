package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const ID_BYTE_LENGTH = 8
const TEXT_LENGTH_BYTE_LENGTH = 2

const DIRNAME = ".data"
const FILENAME = DIRNAME + "/db.dat"

type TextDataRow struct {
	id         uint64
	textLength uint16
	text       string
}

var (
	file *os.File
	err  error
)

func openDB() {
	os.Mkdir(DIRNAME, 0755)

	file, err = os.OpenFile(FILENAME,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0600)
	panic(err)
}

func closeDB() {
	if file != nil {
		for id, cacheRow := range cache {
			if !cacheRow.inMem {
				pushToDisk(id, cacheRow.text)
			}
		}
	} else {
		print("open db first.\n")
	}
	file.Close()
}

func pushToDisk(id uint64, text string) {
	rowToPush := TextDataRow{id, uint16(len(text)), text}
	byteArr := textDataRowToBytes(rowToPush)
	file.Write(byteArr)
}

func getRowFromDisk(id uint64) (TextDataRow, bool) {
	textLengthByteArr := make([]byte, TEXT_LENGTH_BYTE_LENGTH)
	idByteArr := make([]byte, ID_BYTE_LENGTH)
	fileStat, _ := file.Stat()
	size := fileStat.Size()

	for {
		rowId, exit := readId(idByteArr)
		if exit {
			break
		}
		textLength, exit := readTextLength(textLengthByteArr)
		if exit {
			break
		}
		if rowId == id {
			text := readText(textLength)
			return TextDataRow{rowId, textLength, text}, true
		}
		if !seekOverText(textLength, size) {
			break
		}

	}

	return TextDataRow{}, false
}

func readId(idByteArr []byte) (uint64, bool) {
	bytesRead, e := file.Read(idByteArr)
	panic(e)
	if bytesRead == 0 {
		return 0, true
	}
	return binary.LittleEndian.Uint64(idByteArr), false
}

func readTextLength(textLengthByteArr []byte) (uint16, bool) {
	bytesRead, e := file.Read(textLengthByteArr)
	panic(e)
	if bytesRead == 0 {
		return 0, true
	}
	return binary.LittleEndian.Uint16(textLengthByteArr), false
}

func readText(textLength uint16) string {
	textByteArr := make([]byte, int(textLength))
	_, e := file.Read(textByteArr)
	panic(e)
	return string(textByteArr)
}

func seekOverText(textLength uint16, size int64) bool {
	seekLength, e := file.Seek(int64(textLength), 1)
	panic(e)
	return seekLength < size

}

func panic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// for debug purposes
func currentOffSet() {
	b, e := file.Seek(0, 1)
	panic(e)
	fmt.Printf("current byte offset:%d\n", b)
}

/*
 * 64 bits | 16 bit length | var bit length, max length 65536 *note zero length not option
 *    ID   |  textLength   | text
 */
func textDataRowToBytes(row TextDataRow) []byte {
	idByteArr := make([]byte, ID_BYTE_LENGTH)
	textLengthByteArr := make([]byte, TEXT_LENGTH_BYTE_LENGTH)
	var rowByteArr []byte

	textByteArr := []byte(row.text)
	binary.LittleEndian.PutUint64(idByteArr, row.id)
	binary.LittleEndian.PutUint16(textLengthByteArr, row.textLength)

	rowByteArr = append(rowByteArr, idByteArr...)
	rowByteArr = append(rowByteArr, textLengthByteArr...)
	rowByteArr = append(rowByteArr, textByteArr...)

	return rowByteArr
}

func bytesToTextDataRow(b []byte) TextDataRow {
	id := binary.LittleEndian.Uint64(b[:8])
	len := binary.LittleEndian.Uint16(b[8:10])
	text := string(b[10 : len+10])
	return TextDataRow{id, len, text}
}

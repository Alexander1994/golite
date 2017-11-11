package database

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
	osPageSize                    = int32(os.Getpagesize())
	metaTableMaxRowCount          = uint32((osPageSize - metaTableOffsetByteLen) / metaRowByteLen)
	metaTableOffsetAndXtraByteLen = uint32(osPageSize) - (metaTableMaxRowCount * metaRowByteLen)
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

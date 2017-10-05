package main

import "math"

const maxID uint64 = math.MaxUint64 / 2

func idToDBID(id uint64) uint64 {
	if !validID(id) {
		return id
	}
	return (^maxID) | id // add identifier bit
}

func dbIDToID(dbID uint64) uint64 {
	if validID(dbID) {
		return dbID
	}
	return (dbID << 1) >> 1 // pop off identifier bit
}

func validID(ID uint64) bool {
	return ID <= maxID && ID > 0
}

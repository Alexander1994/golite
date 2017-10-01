package main

const UINT_MAX = ^uint64(0)
const MAX_ID = UINT_MAX / 2

func idToDBID(id uint64) uint64 {
	if !validID(id) {
		return id
	}
	return (^MAX_ID) | id // add identifier bit
}

func dbIDToID(dbID uint64) uint64 {
	if validID(dbID) {
		return dbID
	}
	return (dbID << 1) >> 1 // pop off identifier bit
}

func validID(ID uint64) bool {
	return ID <= MAX_ID && ID > 0
}

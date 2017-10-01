package main

import (
	"os"
	"strconv"
)

type CMDEvent func(params []string) // should probably return error/bool

func exitCmd(params []string) { // args should probably be empty
	closeDB()
	os.Exit(0)
}

func insertCmd(params []string) { // %d %s, ID Text
	if len(params) >= 2 {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil && validID(id) {
			insertRow(id, params[1:])
		} else {
			print("first arg should be <=" + string(MAX_ID) + " for the id of the text\n")
		}
	} else {
		print("invalid arg count, 2 expected: %d %s, id text goes here\n")
	}
}

func selectCmd(params []string) { // %d, ID
	if len(params) == 1 {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil && validID(id) {
			findRow(id)
		} else {
			print("could not convert param to int\n")
		}
	} else {
		print("invalid arg count, 1 expected: %d, id \n")
	}
}

func deleteCmd(params []string) { // %d, ID
	if len(params) == 1 {
		if params[0] == "database" {
			resetCache()
			resetPageTable()
			resetDB()
			println("All data in the db and cache has been removed")
		} else {
			id, err := strconv.ParseUint(params[0], 10, 64)
			if err == nil && validID(id) {
				deleteRow(id)
			} else {
				print("could not convert param to int\n")
			}
		}
	} else {
		print("invalid arg count, 1 expected: %d, id \n")
	}
}

func createCmds() map[string]CMDEvent {
	events := make(map[string]CMDEvent)
	events["insert"] = insertCmd
	events["select"] = selectCmd
	events["delete"] = deleteCmd
	events[":exit"] = exitCmd
	return events
}

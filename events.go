
package main

import (
	"os"
	"strconv"
)

type CMDEvent func(params []string) // should probably return error/bool

func exitEvent(params []string) { // args should probably be empty
	closeDB()
	os.Exit(0)
}

func insertEvent(params []string)  { // %d %s, ID Text 
	if (len(params) >= 2 ) {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if (err == nil) {
			insertRow(id, params[1:])
		} else {
			print("first arg should be an uint64 for the id of the text\n")			
		}
	} else {
		print("invalid arg count, 2 expected: %d %s, id text goes here\n")
	}
}

func selectEvent(params []string) { // %d, ID
	if len(params) == 1 {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil {
			findRow(id)
		} else {
			print("could not convert param to int\n")
		}
	} else {
		print("invalid arg count, 1 expected: %d, id \n")
	}
}

func createEvents() map[string]CMDEvent {
	events := make(map[string]CMDEvent)
	events["insert"] = insertEvent
	events["select"] = selectEvent
	events[":exit"] = exitEvent
	return events 
}

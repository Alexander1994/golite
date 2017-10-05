package main

import (
	"fmt"
	"os"
	"strconv"
)

type cmdEvent func(params []string) // should probably return error/bool

func exitCmd(params []string) { // args should probably be empty
	closeDB()
	os.Exit(0)
}

func insertCmd(params []string) { // %d %s, ID Text
	if len(params) >= 2 {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil && validID(id) {
			if !Insert(id, joinOnSpace(params[1:])) {
				print("Unable to insert a row with that ID already exists in the DB\n")
			}

		} else {
			print("first arg should be <=" + string(maxID) + " for the id of the text\n")
		}
	} else {
		print("invalid arg count, 2 expected: %d %s, id text goes here\n")
	}
}

func joinOnSpace(text []string) string {
	textToAddToRow := ""
	for i, param := range text {
		if i != len(text)-1 {
			param += " "
		}
		textToAddToRow += param
	}
	return textToAddToRow
}

func selectCmd(params []string) { // %d, ID
	if len(params) == 1 {
		id, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil && validID(id) {
			text, found := Select(id)
			if found {
				fmt.Printf("%d: %s\n", id, text)
			}
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
			ResetDB()
			println("All data in the db and cache has been removed")
		} else {
			id, err := strconv.ParseUint(params[0], 10, 64)
			if err == nil && validID(id) {
				if Delete(id) {
					print("row deleted from db\n")
				} else {
					print("id not found in memory, row not removed\n")
				}

			} else {
				print("could not convert param to int\n")
			}
		}
	} else {
		print("invalid arg count, 1 expected: %d, id \n")
	}
}

func createCmds() map[string]cmdEvent {
	events := make(map[string]cmdEvent)
	events["insert"] = insertCmd
	events["select"] = selectCmd
	events["delete"] = deleteCmd
	events[":exit"] = exitCmd
	return events
}

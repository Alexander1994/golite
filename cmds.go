package database

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type cmdEvent func(params []string) // should probably return error/bool

func exitCmd(params []string) { // args should probably be empty
	CloseDB()
	os.Exit(0)
}

func createTableCmd(params []string) { // %s tableName
	if len(params) == 1 {
		if strings.Contains(params[0], ".") {
			print("table cannot contain period char")
		} else if !CreateTable(params[0]) {
			print("unable to create table, table already exists\n")
		}
	} else {
		print("invalid arg count, 1 expected: %s, table name goes here\n")
	}
}

func insertCmd(params []string) { // %s %d %s, tableName ID Text
	if len(params) >= 2 {
		rawID, err := strconv.ParseUint(params[1], 10, 32)
		id := uint32(rawID)
		if err == nil && id != 0 {
			if !Insert(id, joinOnSpace(params[2:]), params[0]) {
				print("Unable to insert a row with that ID already exists in the DB\n")
			}

		} else {
			print("second arg should be <=" + strconv.Itoa(math.MaxUint32) + " for the id of the text\n")
		}
	} else {
		print("invalid arg count, 3 expected: %s %d %s ,TableName ID Text\n")
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

func selectCmd(params []string) { // %s %d, tableName ID
	if len(params) == 2 {
		rawID, err := strconv.ParseUint(params[1], 10, 32)
		id := uint32(rawID)
		if err == nil && id != 0 {
			text, found := Select(id, params[0])
			if found {
				fmt.Printf("%d: %s\n", id, text)
			} else {
				print("text not found\n")
			}
		} else {
			print("could not convert param to int\n")
		}
	} else {
		print("invalid arg count, 2 expected: %s %d, TableName ID \n")
	}
}

func deleteCmd(params []string) { // %s %d, tableName ID or %s, tableName
	if len(params) == 2 {
		if params[0] == "database" && params[1] == "confirm" {
			DeleteDB()
			println("All data in the db and cache has been removed")
		} else {
			rawID, err := strconv.ParseUint(params[1], 10, 32)
			id := uint32(rawID)
			if err == nil && id != 0 {
				if Delete(id, params[0]) {
					print("row deleted from db\n")
				} else {
					print("id not found in memory, row not removed\n")
				}

			} else {
				print("could not convert param to int\n")
			}
		}
	} else if len(params) == 1 {
		if !DeleteTable(params[0]) {
			print("table not found, not deleted\n")
		}
	} else {
		print("invalid arg count, 2 expected: %s %d, TableName ID\n")
	}
}

func createCmds() map[string]cmdEvent {
	events := make(map[string]cmdEvent)
	events["insert"] = insertCmd
	events["select"] = selectCmd
	events["delete"] = deleteCmd
	events["create"] = createTableCmd
	events[":exit"] = exitCmd
	return events
}

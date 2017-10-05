package main

/*
// #include <stdlib.h>
import "C"
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func splitOnSpace(c rune) bool {
	return c == ' '
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var cmd string
	cmds := createCmds()
	argsLen := len(os.Args)

	testMode := false
	if argsLen == 2 && os.Args[1] == "test" {
		testMode = true
	}
	openDB(testMode)
	loadPageTable()

	for true {
		print("golite>")
		cmd, _ = reader.ReadString('\n') // second var is error
		cmd = strings.TrimRight(cmd, "\r\n")
		cmdSplitBySpace := strings.FieldsFunc(cmd, splitOnSpace)
		if len(cmdSplitBySpace) > 0 {
			if file == nil {
				openDB(testMode)
			}
			cmdType := cmdSplitBySpace[0]
			event, cmdFound := cmds[cmdType]
			if cmdFound {
				event(cmdSplitBySpace[1:])
			} else {
				fmt.Printf("unrecognized cmd:'%s'\n", cmd)
			}
		} else {
			println("please enter a command")
		}

	}
}

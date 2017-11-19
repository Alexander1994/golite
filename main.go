package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Alexander1994/golite/database"
)

func main() {
	// Setup
	reader := bufio.NewReader(os.Stdin)
	var cmd string
	cmds := createCmds()
	database.OpenDB()

	// Shut down hook
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		args := make([]string, 1)
		cmds[":exit"](args)
	}()

	// Main loop
	for true {
		print("golite>")
		cmd, _ = reader.ReadString('\n') // second var is error
		cmd = strings.TrimRight(cmd, "\r\n")
		cmdSplitBySpace := strings.FieldsFunc(cmd, splitOnSpace)
		if len(cmdSplitBySpace) > 0 {
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

func splitOnSpace(c rune) bool {
	return c == ' '
}

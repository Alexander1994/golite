package database

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var cmd string
	cmds := createCmds()
	OpenDB()

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

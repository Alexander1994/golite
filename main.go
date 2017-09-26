package main
/*
// #include <stdlib.h>
import "C"
*/


import (
	"bufio"
	"os"
	"strings"
	"fmt"
)

func splitOnSpace(c rune) bool {
	return c == ' ' 
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var cmd string
	cmds := createEvents()
	openDB()

	for true {
		print("golite>")
		cmd, _ = reader.ReadString('\n') // second var is error
		cmd = strings.TrimRight(cmd, "\r\n") 
		cmdSplitBySpace := strings.FieldsFunc(cmd, splitOnSpace)
		if (len(cmdSplitBySpace) > 0) {
			cmdType := cmdSplitBySpace[0]
			event, cmdFound := cmds[cmdType]
			if cmdFound {
				event(cmdSplitBySpace[1:])
			} else {
				fmt.Printf("unrecognized cmd:'%s'\n", cmd)			
			}
		} else {
			print("please enter a command\n")
		}			
		
	}
}
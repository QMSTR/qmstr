package main

import (
	"os"
	"strings"

	"github.com/QMSTR/qmstr/pkg/cli"
)

var fixableCommands = map[string]struct{}{"create": struct{}{}, "update": struct{}{}}

func main() {
	os.Args = (*fixCmdLine(&os.Args))
	cli.Execute()
}

func fixCmdLine(cmdline *[]string) *[]string {
	if len(*cmdline) < 3 {
		return cmdline
	}
	if _, ok := fixableCommands[(*cmdline)[1]]; ok {
		nodetype, _, err := cli.TokenizeNodeID((*cmdline)[2])
		if err != nil {
			// invalid cobra will take care
			return cmdline
		}

		if len(*cmdline) > 3 && strings.HasPrefix((*cmdline)[3], nodetype) {
			// nothing to fix
			return cmdline
		}

		// fix command line by inserting subcommand
		newCmdline := append((*cmdline)[:2], append([]string{nodetype}, (*cmdline)[2:]...)...)
		return &newCmdline
	}
	return cmdline
}

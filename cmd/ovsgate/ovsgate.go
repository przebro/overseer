package main

import (
	"os"
	"overseer/cmd/ovsgate/commands"
)

func main() {
	commands.Setup()
	commands.Execute(os.Args)
}

package main

import (
	"os"
	"overseer/cmd/ovswork/commands"
)

func main() {
	commands.Setup()
	commands.Execute(os.Args)
}

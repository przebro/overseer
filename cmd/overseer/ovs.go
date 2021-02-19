package main

import (
	"os"
	"overseer/cmd/overseer/commands"
)

func main() {
	commands.Setup()
	commands.Execute(os.Args)
}

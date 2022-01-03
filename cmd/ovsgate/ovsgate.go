package main

import (
	"os"

	"github.com/przebro/overseer/cmd/ovsgate/commands"
)

func main() {
	commands.Setup()
	commands.Execute(os.Args)
}

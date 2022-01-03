package main

import (
	"strings"

	"github.com/przebro/overseer/cmd/ovstools/commands"
	"github.com/przebro/overseer/ovscli"

	"github.com/manifoldco/promptui"
)

func main() {

	commands.Setup(ovscli.CreateClient())

	for {

		prompt := promptui.Prompt{Label: "ovscli",
			Templates: &promptui.PromptTemplates{
				Prompt: ``,
			},
		}
		args, _ := prompt.Run()
		commands.Execute(strings.Split(args, " "))
	}

}

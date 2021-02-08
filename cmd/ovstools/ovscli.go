package main

import (
	"overseer/cmd/ovstools/commands"
	"overseer/ovscli"
	"strings"

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

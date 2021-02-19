package commands

import (
	"overseer/ovscli"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ovscli",
	Short: "command line tools for overseer",
}

var ovsAddr string

//Setup - performs setup
func Setup(client *ovscli.OverseerClient) {

	rootCmd.AddCommand(createQuitCommand(client))
	rootCmd.AddCommand(createSetsrvCommand(client))
	rootCmd.AddCommand(createLoginCommand(client))
	rootCmd.AddCommand(createAddCmd(client))
	rootCmd.AddCommand(createDelCmd(client))
	rootCmd.AddCommand(createCheckCmd(client))
	rootCmd.AddCommand(createListCmd(client))
	rootCmd.AddCommand(createSetCmd(client))
	rootCmd.AddCommand(createRemoveCmd(client))
	rootCmd.AddCommand(createOrderCmd(client))
	rootCmd.AddCommand(createTaskCmd(client))
	rootCmd.AddCommand(createPoolCmd(client))
}

//Execute - executes commands
func Execute(args []string) error {

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

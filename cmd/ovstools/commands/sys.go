package commands

import (
	"fmt"
	"os"
	"overseer/ovscli"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var login = &cobra.Command{}
var isSecure bool

func createSetsrvCommand(client *ovscli.OverseerClient) *cobra.Command {

	var certpath string
	cmd := &cobra.Command{
		Use:     "setsrv",
		Short:   "setup connection to server",
		Example: "setsrv localhost:7053 --cert=path_to_cert_file",
		Args:    cobra.ExactArgs(1),
		Run:     func(c *cobra.Command, args []string) { setupConnection(client, c, args, certpath) },
	}
	cmd.Flags().StringVar(&certpath, "cert", "", "path to certifictate")

	return cmd
}

func createQuitCommand(client *ovscli.OverseerClient) *cobra.Command {

	quit := &cobra.Command{
		Use:   "quit",
		Short: "quits from application",
		Run:   func(c *cobra.Command, args []string) { runQuit(client, c, args) },
	}
	return quit

}

func createLoginCommand(client *ovscli.OverseerClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Run:   func(c *cobra.Command, args []string) { loginCommand(client, c, args) },
		Short: "authenticate user against server with username and password",
	}

	return cmd
}

func setupConnection(client *ovscli.OverseerClient, cmd *cobra.Command, args []string, certpath string) {

	result := client.Connect(args[0], certpath)
	fmt.Println(result)

}

func runQuit(client *ovscli.OverseerClient, cmd *cobra.Command, args []string) {

	client.Close()
	os.Exit(0)
}

func loginCommand(client *ovscli.OverseerClient, cmd *cobra.Command, args []string) {

	var username string
	var passwd string
	var token string
	var err error
	uprompt := promptui.Prompt{Label: "username"}

	if username, err = uprompt.Run(); err != nil {
		fmt.Println(err)
		return
	}

	pprompt := promptui.Prompt{Label: "password", Mask: '*'}

	if passwd, err = pprompt.Run(); err != nil {
		fmt.Println(err)
		return
	}

	if token, err = client.Authenticate(username, passwd); err != nil {
		fmt.Println(err)
	}

	fmt.Println(token)

}

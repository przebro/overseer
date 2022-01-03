package commands

import (
	"fmt"
	"os"
	"overseer/common/types"
	"overseer/ovscli"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func createSetsrvCommand(client *ovscli.OverseerClient) *cobra.Command {

	var serverCA string
	var certificatePath string
	var certificateKeyPath string

	cmd := &cobra.Command{
		Use:     "setsrv",
		Short:   "setup connection to server",
		Example: "setsrv localhost:7053 --rootca=path_to_cert_file --cert=path_to_cert --key=path_to_key",
		Args:    cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			setupConnection(client, c, args, serverCA, certificatePath, certificateKeyPath)
		},
	}
	cmd.Flags().StringVar(&serverCA, "rootca", "", "path to server's rootCA")
	cmd.Flags().StringVar(&certificatePath, "cert", "", "path to certifictate")
	cmd.Flags().StringVar(&certificateKeyPath, "key", "", "path to client's key")

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

func setupConnection(client *ovscli.OverseerClient, cmd *cobra.Command, args []string, serverCA, clientCertPath, clientKeyPath string) {

	fmt.Println(serverCA)

	var level types.ConnectionSecurityLevel = types.ConnectionSecurityLevelNone
	var policy types.CertPolicy = types.CertPolicyNone

	if serverCA != "" {
		level = types.ConnectionSecurityLevelServeOnly
		policy = types.CertPolicyRequired
	}

	if clientCertPath != "" && clientKeyPath != "" {
		level = types.ConnectionSecurityLevelClientAndServer
	}

	result := client.Connect(args[0], serverCA, clientCertPath, clientKeyPath, level, policy)
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

package commands

import (
	"fmt"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/ovscli"

	"github.com/spf13/cobra"
)

var validStsates = map[string]int32{"SHR": 0, "EXL": 1}

func createAddCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "ADD",
		Short: "ADD - adds a ticket",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {

			n, o := unfold(args)
			clientAddTicket(client, c, n, o)

		},
	}
	return cmd
}

func createDelCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "DEL",
		Short: "DEL - deletes a ticket",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {

			n, o := unfold(args)
			clientRemoveTicket(client, c, n, o)

		},
	}
	return cmd
}

func createCheckCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "CHECK",
		Short: "CHECK - Checks if  a ticket exists",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {

			n, o := unfold(args)
			clientCheckTicket(client, c, n, o)

		},
	}
	return cmd
}

func createListCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "TICKETS",
		Short: "TICKETS - List tickets",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {

			n, o := unfold(args)
			clientListTickets(client, c, n, o)
		},
	}
	return cmd
}

func createSetCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "SET",
		Short: "SET - sets a flag",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {

			n, s := unfold(args)
			clientSetFlag(client, c, n, s)
		},
	}
	return cmd
}

func createRemoveCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "REM",
		Short: "REM -  removes a flag permanently",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {

			clientRemoveFlag(client, c, args[0])
		},
	}
	return cmd
}

func clientAddTicket(client *ovscli.OverseerClient, cmd *cobra.Command, name, odate string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(date.Odate(odate), "odate"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.AddTicket(name, odate)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println(result)

}

func clientRemoveTicket(client *ovscli.OverseerClient, cmd *cobra.Command, name, odate string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(date.Odate(odate), "odate"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.DelTicket(name, odate)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}

func clientCheckTicket(client *ovscli.OverseerClient, cmd *cobra.Command, name, odate string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(date.Odate(odate), "odate"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.CheckTicket(name, odate)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func clientListTickets(client *ovscli.OverseerClient, cmd *cobra.Command, name, odate string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(date.Odate(odate), "odate"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.ListTickets(name, odate)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, x := range result {
		fmt.Printf("Name:%s ODAT:%s\n", x.Name, x.Odate)
	}
}

func clientSetFlag(client *ovscli.OverseerClient, cmd *cobra.Command, name, state string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	s, ok := validStsates[state]

	if !ok {
		fmt.Printf("unknown flag state:%s, valid states are: SHR | EXL", state)
		return
	}

	result, err := client.SetFlag(name, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func clientRemoveFlag(client *ovscli.OverseerClient, cmd *cobra.Command, name string) {

	if err := validator.Valid.ValidateTag(name, "resvalue,max=32"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.DestroyFlag(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func unfold(args []string) (string, string) {
	var name string
	var val string

	if len(args) == 2 {
		name, val = args[0], args[1]
	}
	if len(args) == 1 {
		name = args[0]
	}

	return name, val
}

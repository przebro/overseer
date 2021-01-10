package main

import (
	"errors"
	"flag"
	"fmt"
	"overseer/ovscli"
	"os"
)

var orderCmdSet *flag.FlagSet
var holdCmdSet *flag.FlagSet
var listCmdSet *flag.FlagSet
var infoCmdSet *flag.FlagSet
var okCmdSet *flag.FlagSet
var rerunCmdSet *flag.FlagSet

var resourceName string
var resourceGroup string
var ticketOdate string
var forceTask bool

const envNameOvs string = "OVSTOOLSRV"

func init() {

	orderCmdSet = flag.NewFlagSet("ORDER", flag.ExitOnError)
	orderCmdSet.StringVar(&resourceName, "n", "", "Definition name")
	orderCmdSet.StringVar(&resourceGroup, "g", "", "Definition group")
	orderCmdSet.StringVar(&ticketOdate, "d", "", "Odate value")
	orderCmdSet.BoolVar(&forceTask, "f", false, "Force task")

	listCmdSet = flag.NewFlagSet("LIST", flag.ExitOnError)

	infoCmdSet = flag.NewFlagSet("INFO", flag.ExitOnError)
	infoCmdSet.StringVar(&resourceName, "o", "", "OrderID")

	okCmdSet = flag.NewFlagSet("SETOK", flag.ExitOnError)
	okCmdSet.StringVar(&resourceName, "o", "", "OrderID")

	rerunCmdSet = flag.NewFlagSet("RERUN", flag.ExitOnError)
	rerunCmdSet.StringVar(&resourceName, "o", "", "OrderID")

}

func main() {

	var err error
	var cli *ovscli.OverseerClient
	var resultMsg string

	if len(os.Args) < 2 {
		fmt.Println("Wrong number of arguments")
		printDefaults()
		os.Exit(8)
	}

	switch os.Args[1] {
	case "ORDER":
		{
			orderCmdSet.Parse(os.Args[2:])
		}
	case "LIST":
		{
			listCmdSet.Parse(os.Args[2:])
		}
	case "INFO":
		{
			infoCmdSet.Parse(os.Args[2:])
		}
	case "SETOK":
		{
			okCmdSet.Parse(os.Args[2:])
		}
	case "RERUN":
		{
			rerunCmdSet.Parse(os.Args[2:])
		}
	default:
		{
			fmt.Println("unrecognized option")
			printDefaults()
			os.Exit(8)
		}
	}

	addrValue := os.Getenv(envNameOvs)
	if addrValue == "" {
		fmt.Println(fmt.Sprintf("Env variable %s not found. Set this variable to address of ovs server.", envNameOvs))
		os.Exit(8)
	}

	if orderCmdSet.Parsed() {
		err = validateName(resourceName)
		if err != nil {

		}
		err = validateOdate(ticketOdate)
		if err != nil {

		}

		cli, err = ovscli.NewClient(addrValue)
		err = validateName(resourceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}

		resultMsg = cli.OrderTask(resourceGroup, resourceName, ticketOdate, forceTask)
		fmt.Println(resultMsg)

	}
	if listCmdSet.Parsed() {
		fmt.Println("Listing tasks")
		cli, err = ovscli.NewClient(addrValue)
		result := cli.ListTasks()
		writeList(result)

	}
	if infoCmdSet.Parsed() {
		cli, err = ovscli.NewClient(addrValue)
		result := cli.TaskDetail(resourceName)
		for r := range result {
			fmt.Println(r)
		}
	}
	if okCmdSet.Parsed() {
		fmt.Println("Setting task to Ended OK")
		cli, err = ovscli.NewClient(addrValue)
		result := cli.SetToOK(resourceName)
		fmt.Println(result)

	}
	if rerunCmdSet.Parsed() {
		fmt.Println("Setting task to Ended OK")
		cli, err = ovscli.NewClient(addrValue)
		result := cli.Rerun(resourceName)
		fmt.Println(result)

	}

}

func validateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if name != "" && len(name) > 20 {
		return errors.New("name length cannot be greater than 20")
	}

	return nil
}
func validateOdate(odate string) error {
	if odate != "" && len(odate) != 8 {
		return errors.New("invalid odate length")
	}

	return nil
}

func writeList(data []ovscli.ListResult) {

	fmt.Println("----------------------------------------------------")
	for x, e := range data {
		fmt.Printf("%d %s %15s  %15s %10s\n", x, e.ID, e.Name, e.Group, e.Status)
		for _, n := range e.Info {
			fmt.Println(n)
		}
	}

}

func printDefaults() {
	fmt.Println("ovstask ORDER - Order a new task")
	orderCmdSet.PrintDefaults()
	fmt.Println("ovstask LIST - list active tasks")
	listCmdSet.PrintDefaults()
	fmt.Println("ovstask INFO - get task detail")
	infoCmdSet.PrintDefaults()
	fmt.Println("ovstask SETOK - set a task to Ended OK")
	okCmdSet.PrintDefaults()
	fmt.Println("ovstask RERUN - rerun a task")
	rerunCmdSet.PrintDefaults()
}

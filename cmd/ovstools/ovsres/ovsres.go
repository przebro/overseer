package main

import (
	"errors"
	"flag"
	"fmt"
	"overseer/ovscli"
	"os"
)

var addCmdSet *flag.FlagSet
var checkCmdSet *flag.FlagSet
var delCmdSet *flag.FlagSet
var setCmdSet *flag.FlagSet

var resourceName string
var ticketOdate string
var flagValue string
var flagSwitch bool
var ticketSwitch bool

var flagState = map[string]int32{
	"NONE": 0,
	"SHR":  1,
	"EXL":  2,
}

const envNameOvs string = "OVSTOOLSRV"

func init() {
	addCmdSet = flag.NewFlagSet("ADD", flag.ExitOnError)
	addCmdSet.StringVar(&resourceName, "n", "", "Resource name")
	addCmdSet.StringVar(&ticketOdate, "d", "", "Odate value")

	checkCmdSet = flag.NewFlagSet("CHECK", flag.ExitOnError)
	checkCmdSet.StringVar(&resourceName, "n", "", "Resource name")
	checkCmdSet.BoolVar(&flagSwitch, "f", false, "")
	checkCmdSet.BoolVar(&ticketSwitch, "t", false, "")
	checkCmdSet.StringVar(&flagValue, "v", "", "Flag value [EXL|SHR|NONE]")
	checkCmdSet.StringVar(&ticketOdate, "d", "", "Odate value")

	delCmdSet = flag.NewFlagSet("DEL", flag.ExitOnError)
	delCmdSet.StringVar(&resourceName, "n", "", "Resource name")
	delCmdSet.StringVar(&ticketOdate, "d", "", "Odate value")

	setCmdSet = flag.NewFlagSet("SET", flag.ExitOnError)
	setCmdSet.StringVar(&resourceName, "n", "", "Resource name")
	setCmdSet.StringVar(&flagValue, "v", "", "Flag value [EXL|SHR|NONE]")

}

func main() {

	var err error
	var cli *ovscli.OverseerClient
	var resultMsg string

	if len(os.Args) < 2 {
		fmt.Println("Wrong number of arguments")
		fmt.Println("ovsres ADD - Add new ticket")
		addCmdSet.PrintDefaults()
		fmt.Println("ovsres CHECK - Returns information if ticket or flag exists")
		checkCmdSet.PrintDefaults()
		fmt.Println("ovsres DEL - Deletes ticket")
		delCmdSet.PrintDefaults()
		fmt.Println("ovsres SET - Sets flag")
		setCmdSet.PrintDefaults()

		os.Exit(8)
	}

	switch os.Args[1] {
	case "ADD":
		{
			err = addCmdSet.Parse(os.Args[2:])
		}
	case "DEL":
		{
			err = delCmdSet.Parse(os.Args[2:])
		}
	case "CHECK":
		{
			err = checkCmdSet.Parse(os.Args[2:])
		}
	case "SET":
		{
			err = setCmdSet.Parse(os.Args[2:])
		}
	default:
		{
			fmt.Println("Unrecognized option.")
			flag.PrintDefaults()
			os.Exit(8)
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	addrValue := os.Getenv(envNameOvs)
	if addrValue == "" {
		fmt.Println(fmt.Sprintf("Env variable %s not found. Set this variable to address of ovs server.", envNameOvs))
		os.Exit(8)
	}

	err = validateName(resourceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	if setCmdSet.Parsed() {

		val, e := flagState[flagValue]
		if e == false {
			fmt.Println("Invalid flag state posible values are: NONE | EXL | SHR")
			os.Exit(8)
		}

		cli, err = ovscli.NewClient(addrValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		resultMsg = cli.SetFlag(resourceName, val)

	}
	if checkCmdSet.Parsed() {

		if flagSwitch == ticketSwitch {
			fmt.Println("only one flag f | t can be specified.")
			os.Exit(8)
		}

		cli, err = ovscli.NewClient(addrValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		if flagSwitch {
			val, e := flagState[flagValue]
			if e == false {
				fmt.Println("Invalid flag state posible values are: NONE | EXL | SHR")
				os.Exit(8)
			}
			resultMsg = cli.CheckFlag(resourceName, val)
		} else {
			resultMsg = cli.CheckTicket(resourceName, ticketOdate)
		}

	}
	if addCmdSet.Parsed() || delCmdSet.Parsed() {

		err = validateOdate(ticketOdate)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		cli, err = ovscli.NewClient(addrValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		if addCmdSet.Parsed() {
			resultMsg = cli.AddTicket(resourceName, ticketOdate)
		} else {
			resultMsg = cli.DelTicket(resourceName, ticketOdate)
		}
	}

	fmt.Println(resultMsg)

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

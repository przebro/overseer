package commands

import (
	"fmt"
	"strings"

	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/ovscli"

	"github.com/spf13/cobra"
)

type taskAction int

const (
	taskActionSetOk   taskAction = 1
	taskActionHold    taskAction = 2
	taskActionFree    taskAction = 3
	taskActionConfirm taskAction = 4
	taskActionRerun   taskAction = 5
	taskActionShow    taskAction = 6
)

func createOrderCmd(client *ovscli.OverseerClient) *cobra.Command {

	var force bool

	cmd := &cobra.Command{

		Use:   "ORDER",
		Short: "ORDER - Puts a definition into active task pool",
		Args:  cobra.ExactArgs(3),
		Run: func(c *cobra.Command, args []string) {

			clientOrderDefinition(client, args[0], args[1], args[2], force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "forcefully order task")

	return cmd
}

func createTaskCmd(client *ovscli.OverseerClient) *cobra.Command {

	var setok bool
	var hold bool
	var free bool
	var confirm bool
	var rerun bool
	var show bool

	var action taskAction

	cmd := &cobra.Command{

		Use:   "TASK",
		Short: "TASK - Perofrms an anction on a task",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {

			var selected int
			var ok bool

			if setok == false && hold == false && free == false && confirm == false && rerun == false && show == false {
				return fmt.Errorf("a flag is required")
			}
			acts := map[int]taskAction{
				1:  taskActionSetOk,
				2:  taskActionHold,
				4:  taskActionFree,
				8:  taskActionConfirm,
				16: taskActionRerun,
				32: taskActionShow,
			}

			if setok {
				selected++
				setok = false
			}

			if hold {
				selected += 2
				hold = false
			}

			if free {
				selected += 4
				free = false
			}

			if confirm {
				selected += 8
				confirm = false
			}

			if rerun {
				selected += 16
				rerun = false
			}

			if show {
				selected += 32
				show = false
			}
			fmt.Println("selected in check:", selected)

			if action, ok = acts[selected]; ok == false {
				return fmt.Errorf("one and only one option can be selected")
			}

			return nil
		},
		Run: func(c *cobra.Command, args []string) {

			orderID := args[0]
			clientTaskAction(client, orderID, action)
		},
	}
	cmd.Flags().BoolVar(&setok, "setok", false, "sets task in an ended ok status")
	cmd.Flags().BoolVar(&hold, "hold", false, "holds procesing task")
	cmd.Flags().BoolVar(&free, "free", false, "frees holded task")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "manually confirms a task")
	cmd.Flags().BoolVar(&rerun, "rerun", false, "reruns a task")
	cmd.Flags().BoolVar(&show, "show", false, "shows details of a task")

	return cmd
}

func createPoolCmd(client *ovscli.OverseerClient) *cobra.Command {

	cmd := &cobra.Command{

		Use:   "POOL",
		Short: "POOL - shows al list of active tasks in a task pool",
		Run: func(c *cobra.Command, args []string) {
			clientShowTaskpool(client)
		},
	}

	return cmd
}

func clientOrderDefinition(client *ovscli.OverseerClient, group, name, odate string, force bool) {

	if err := validator.Valid.ValidateTag(group, "required,max=20,resname"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(name, "required,max=32,resname"); err != nil {
		fmt.Println(err)
		return
	}

	if err := validator.Valid.ValidateTag(name, "required,odate"); err != nil {
		fmt.Println(err)
		return
	}

	client.OrderTask(group, name, odate, force)
}

func clientTaskAction(client *ovscli.OverseerClient, orderID string, action taskAction) {

	var result string
	var err error

	if err := validator.Valid.ValidateTag(orderID, "required,len=5"); err != nil {
		fmt.Println(err)
		return
	}

	if action == taskActionSetOk {
		result, err = client.SetToOK(orderID)
	}

	if action == taskActionHold {
		result, err = client.Hold(orderID)
	}
	if action == taskActionFree {
		result, err = client.Free(orderID)
	}

	if action == taskActionConfirm {
		result, err = client.Confirm(orderID)
	}

	if action == taskActionRerun {
		result, err = client.Rerun(orderID)
	}

	if action == taskActionShow {
		result, err = client.TaskDetail(orderID)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	r := strings.Split(result, "\n")
	for _, n := range r {
		fmt.Println(n)
	}

}

func clientShowTaskpool(client *ovscli.OverseerClient) {

	var result []ovscli.ListResult
	var err error
	if result, err = client.ListTasks(); err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		fmt.Printf("%6s %20s %20s %s\n", v.ID, v.Group, v.Name, v.Status)
		for _, n := range v.Info {
			fmt.Println(n)
		}
	}
}

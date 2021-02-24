package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("OVS - check program")
	fmt.Println("--------------------------------------------------")
	fmt.Println("Arguments")
	fmt.Println("--------------------------------------------------")
	for x, e := range os.Args {
		fmt.Println(x, ":", e)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("Variables")
	fmt.Println("--------------------------------------------------")

	env := os.Environ()
	vars := make([]string, 0)
	for x := range env {
		if strings.HasPrefix(env[x], "OVS_") {
			vars = append(vars, env[x])
		}
	}

	for x, e := range vars {
		fmt.Println(x, ":", e)

	}

	timeOut := 0

	abend := os.Getenv("OVS_ABEND")
	timeout := os.Getenv("OVS_TIMEOUT")
	rn := os.Getenv("OVS_RN")

	if timeout != "" {
		timeOut, _ = strconv.Atoi(timeout)

	}

	if timeOut > 0 {
		for i := 0; i < timeOut; i++ {
			fmt.Println("doing sleep until timeout...")
			time.Sleep(1 * time.Second)
		}
	}

	if abend == "Y" && rn == "1" {
		os.Exit(5)
	}

	os.Exit(0)

}

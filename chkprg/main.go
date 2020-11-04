package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	runNumber := 0

	abend := os.Getenv("OVS_ABEND")
	rn := os.Getenv("OVS_RN")
	if rn != "" {
		runNumber, _ = strconv.Atoi(rn)

	}

	if abend == "Y" && runNumber == 1 {
		os.Exit(5)
	}
	os.Exit(0)

}

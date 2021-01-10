package main

import (
	"flag"
	"fmt"
	"overseer/common/helpers"
	"overseer/common/logger"
	"overseer/ovsworker"
	"overseer/ovsworker/config"
	"os"
	"path/filepath"
)

var (
	conf     *config.Config
	wprofile bool
)

func init() {
	flag.BoolVar(&wprofile, "wprofile", false, "Start profiler")

}

func main() {

	root, _, err := helpers.GetDirectories(os.Args[0])

	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	flag.Parse()
	if !flag.Parsed() {
		fmt.Println("unable to parse flags")
		flag.PrintDefaults()
		os.Exit(8)
	}

	conf, err := config.Load(filepath.Join(root, "config", "worker.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	logcfg := conf.GetLogConfiguration()
	log := logger.NewLogger(logcfg.Directory, logcfg.Level)

	worker := ovsworker.NewWorkerService(conf)

	if wprofile == true {
		helpers.StartProfiler(log, "workerprofile.prof")
	}

	err = worker.Start()
	if err != nil {
		log.Error(err)
		os.Exit(8)
	}
	os.Exit(0)

}

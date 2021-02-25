package main

import (
	"flag"
	"fmt"
	"os"
	"overseer/common/helpers"
	"overseer/common/logger"
	"overseer/ovsworker"
	"overseer/ovsworker/config"
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

	var rootPath string
	var progPath string
	var err error
	if rootPath, progPath, err = helpers.GetDirectories(os.Args[0]); err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	flag.Parse()
	if !flag.Parsed() {
		fmt.Println("unable to parse flags")
		flag.PrintDefaults()
		os.Exit(8)
	}

	conf, err := config.Load(filepath.Join(rootPath, "config", "worker.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	logcfg := conf.GetLogConfiguration()
	log := logger.NewLogger(logcfg.Directory, logcfg.Level)

	conf.Worker.ProcessDirectory = progPath
	conf.Worker.RootDirectory = rootPath

	worker := ovsworker.NewWorkerService(conf)

	if worker == nil {
		os.Exit(8)
	}

	if wprofile == true {
		helpers.StartProfiler(log, "workerprofile.prof")
	}

	err = worker.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}
	os.Exit(0)

}

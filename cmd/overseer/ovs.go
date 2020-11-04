package main

import (
	"flag"
	"fmt"
	"goscheduler/common/helpers"
	"goscheduler/common/logger"
	"goscheduler/overseer"
	"goscheduler/overseer/config"
	"os"
	"path/filepath"
)

const (
	configurationDirectory = "config"
)

var (
	conf       *config.OverseerConfiguration
	configFile string
	hostAddr   string
	hostPort   int
	profile    bool
)

func init() {

	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.StringVar(&hostAddr, "host", "", "Host address")
	flag.IntVar(&hostPort, "port", 0, "Host port")
	flag.BoolVar(&profile, "profile", false, "enable profiling")

}

func main() {

	root, prog, err := helpers.GetDirectories(os.Args[0])
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

	config, err := getConfiguration(root, prog)
	//Get log section from configuration
	logcfg := config.GetLogConfiguration()
	log := logger.NewLogger(logcfg.LogDirectory, logcfg.LogLevel)

	if err != nil {
		log.Error("failed to read configuration")
		os.Exit(8)
	}

	ovs, err := overseer.NewInstance(*config)
	if err != nil {
		log.Error(err)
	}

	if profile == true {
		helpers.StartProfiler(log, "schedulerprofile.prof")
	}

	err = ovs.Start()

	if err != nil {
		log.Error(err)
		os.Exit(8)
	}

	os.Exit(0)

}

func getConfiguration(root, prog string) (*config.OverseerConfiguration, error) {

	var err error
	//If flag is specified, check the custom configuration file.
	if configFile != "" {

		if conf, err = config.Load(configFile); err != nil {
			return nil, err
		}

	} else {
		//Use a built-in default configuration
		if conf, err = config.Load(filepath.Join(root, "config", "overseer.json")); err != nil {
			return nil, err
		}
	}

	conf.ProcessDirectory = prog
	conf.RootDirectory = root

	return conf, nil

}

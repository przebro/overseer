package commands

import (
	"fmt"
	"os"
	"overseer/common/core"
	"overseer/common/helpers"
	"overseer/common/logger"
	"overseer/common/validator"
	"overseer/ovsworker"
	"overseer/ovsworker/config"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	configFileName = "worker.json"
	configDirName  = "config"
	profFileName   = "workerprofile.prof"
)

var (
	configFile string
	hostAddr   string
	hostPort   int
	profile    bool
)

var rootCmd = &cobra.Command{
	Use:   "ovs",
	Short: "overseer scheduler",
	Run: func(c *cobra.Command, args []string) {
		startWorker()
	},
}

//Setup - performs setup
func Setup() {
	rootCmd.Flags().StringVar(&configFile, "config", "", "path to configuration file")
	rootCmd.Flags().StringVar(&hostAddr, "host", "", "overseer address")
	rootCmd.Flags().IntVar(&hostPort, "port", 0, "overseer port")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "starts profiling")
}

//Execute - executes commands
func Execute(args []string) error {

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func startWorker() {

	var rootPath string
	var progPath string
	var err error
	var worker core.RunnableComponent
	var conf config.OverseerWorkerConfiguration

	if rootPath, progPath, err = helpers.GetDirectories(os.Args[0]); err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	if conf, err = getConfiguration(rootPath, progPath); err != nil {
		fmt.Println(err)
		os.Exit(16)
	}

	logcfg := conf.GetLogConfiguration()
	log := logger.NewLogger(logcfg.Directory, logcfg.Level)

	if worker, err = ovsworker.NewWorkerService(conf); err != nil {
		log.Error(err)
		os.Exit(8)
	}

	log.Info("starting runner")

	if profile == true {
		helpers.StartProfiler(log, profFileName)
	}

	runner := core.NewRunner(worker)
	runner.Run()

	if err != nil {
		log.Error(err)
		os.Exit(8)
	}

	os.Exit(0)

}

func getConfiguration(root, prog string) (config.OverseerWorkerConfiguration, error) {

	var err error
	var conf config.OverseerWorkerConfiguration
	//If flag is specified, check the custom configuration file.
	if configFile != "" {

		if conf, err = config.Load(configFile); err != nil {
			return conf, err
		}

	} else {
		//Use a built-in default configuration
		if conf, err = config.Load(filepath.Join(root, configDirName, configFileName)); err != nil {
			return conf, err
		}
		if err = validator.Valid.Validate(conf); err != nil {
			return conf, err
		}
	}

	if hostAddr != "" {
		conf.Worker.Host = hostAddr
	}

	if hostPort != 0 {
		conf.Worker.Port = hostPort
	}

	conf.Worker.ProcessDirectory = prog
	conf.Worker.RootDirectory = root

	if err = validator.Valid.Validate(conf); err != nil {
		return conf, err
	}

	return conf, nil

}

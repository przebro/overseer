package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/common/helpers"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/ovsworker"
	"github.com/przebro/overseer/ovsworker/config"
	"github.com/rs/zerolog/log"

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

// Setup - performs setup
func Setup() {
	rootCmd.Flags().StringVar(&configFile, "config", "", "path to configuration file")
	rootCmd.Flags().StringVar(&hostAddr, "host", "", "overseer address")
	rootCmd.Flags().IntVar(&hostPort, "port", 0, "overseer port")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "starts profiling")
}

// Execute - executes commands
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

	if !filepath.IsAbs(logcfg.LogDirectory) {
		logcfg.LogDirectory = filepath.Join(rootPath, logcfg.LogDirectory)
	}

	if _, err := os.Stat(logcfg.LogDirectory); os.IsNotExist(err) {
		os.MkdirAll(logcfg.LogDirectory, 0755)
	}

	logger.Configure("ovs_worker", logcfg)

	if !filepath.IsAbs(conf.Worker.SysoutDirectory) {
		conf.Worker.SysoutDirectory = filepath.Join(rootPath, conf.Worker.SysoutDirectory)
	}

	if _, err := os.Stat(conf.Worker.SysoutDirectory); os.IsNotExist(err) {
		os.MkdirAll(conf.Worker.SysoutDirectory, 0755)
	}

	if worker, err = ovsworker.NewWorkerService(conf); err != nil {
		log.Error().Err(err).Msg("error creating worker")
		os.Exit(8)
	}

	log.Info().Msg("starting runner")

	if profile == true {
		helpers.StartProfiler(log.Logger, profFileName)
	}

	runner := core.NewRunner(worker)
	runner.Run()

	if err != nil {
		log.Error().Err(err).Msg("error starting worker")
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

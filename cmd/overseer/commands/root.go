package commands

import (
	"fmt"
	"os"
	"overseer/common/core"
	"overseer/common/helpers"
	"overseer/common/logger"
	"overseer/common/validator"
	"overseer/overseer"
	"overseer/overseer/config"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	configFileName = "overseer.json"
	configDirName  = "config"
	profFileName   = "schedulerprofile.prof"
)

var (
	configFile string
	hostAddr   string
	hostPort   int
	profile    bool
	quiesce    bool
)

var rootCmd = &cobra.Command{
	Use:   "ovs",
	Short: "overseer scheduler",
	Run: func(c *cobra.Command, args []string) {
		startOverseer()
	},
}

//Setup - performs setup
func Setup() {
	rootCmd.Flags().StringVar(&configFile, "config", "", "path to configuration file")
	rootCmd.Flags().StringVar(&hostAddr, "host", "", "overseer address")
	rootCmd.Flags().IntVar(&hostPort, "port", 0, "overseer port")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "starts profiling")
	rootCmd.Flags().BoolVar(&quiesce, "quiesce", false, "starts overseer as service")
}

//Execute - executes commands
func Execute(args []string) error {

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func startOverseer() {

	var rootPath string
	var progPath string
	var err error
	var ovs core.RunnableComponent
	var conf config.OverseerConfiguration
	var log logger.AppLogger

	if rootPath, progPath, err = helpers.GetDirectories(os.Args[0]); err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	if conf, err = getConfiguration(rootPath, progPath); err != nil {
		fmt.Println(err)
		os.Exit(16)
	}

	//Get log section from configuration
	logcfg := conf.GetLogConfiguration()

	if !filepath.IsAbs(logcfg.LogDirectory) {
		logcfg.LogDirectory = filepath.Join(rootPath, logcfg.LogDirectory)
	}

	if log, err = logger.NewLogger(logcfg); err != nil {
		fmt.Println(err)
		os.Exit(16)
	}

	if ovs, err = overseer.New(conf, log, quiesce); err != nil {
		log.Error(err)
		os.Exit(16)
	}

	if profile == true {
		helpers.StartProfiler(log, profFileName)
	}

	runner := core.NewRunner(ovs)
	err = runner.Run()

	if err != nil {
		log.Error(err)
		os.Exit(8)
	}

	os.Exit(0)

}

func getConfiguration(root, prog string) (config.OverseerConfiguration, error) {

	var err error
	var conf config.OverseerConfiguration
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
	}

	if hostAddr != "" {
		conf.Server.Host = hostAddr
	}

	if hostPort != 0 {
		conf.Server.Port = hostPort
	}

	conf.Server.ProcessDirectory = prog
	conf.Server.RootDirectory = root

	if err = validator.Valid.Validate(conf); err != nil {
		return conf, err
	}

	return conf, nil

}

package commands

import (
	"fmt"
	"os"
	"overseer/common/helpers"
	"overseer/common/logger"
	"overseer/overseer"
	"overseer/overseer/config"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	conf       *config.OverseerConfiguration
	configFile string
	hostAddr   string
	hostPort   int
	profile    bool
)

var rootCmd = &cobra.Command{
	Use:   "ovs",
	Short: "command line tools for overseer",
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
	var ovs *overseer.Overseer

	if rootPath, progPath, err = helpers.GetDirectories(os.Args[0]); err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	if err = getConfiguration(rootPath, progPath); err != nil {
		fmt.Println(err)
		os.Exit(16)
	}

	//Get log section from configuration
	logcfg := conf.GetLogConfiguration()
	log := logger.NewLogger(logcfg.LogDirectory, logcfg.LogLevel)

	if ovs, err = overseer.NewInstance(*conf); err != nil {
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

func getConfiguration(root, prog string) error {

	var err error
	//If flag is specified, check the custom configuration file.
	if configFile != "" {

		if conf, err = config.Load(configFile); err != nil {
			return err
		}

	} else {
		//Use a built-in default configuration
		if conf, err = config.Load(filepath.Join(root, "config", "overseer.json")); err != nil {
			return err
		}
	}

	conf.Server.ProcessDirectory = prog
	conf.Server.RootDirectory = root

	return nil

}

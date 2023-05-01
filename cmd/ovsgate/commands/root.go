package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/przebro/overseer/common/helpers"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/ovsgate"
	"github.com/przebro/overseer/ovsgate/config"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	configFile string
	hostAddr   string
	hostPort   int
	profile    bool
)

var rootCmd = &cobra.Command{
	Run: func(c *cobra.Command, args []string) {
		startGateway()
	},
}

// Setup - initializes root command
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

func startGateway() {

	var rootPath string
	var progPath string
	var err error
	var gate *ovsgate.OverseerGateway
	var conf config.OverseerGatewayConfig

	if rootPath, progPath, err = helpers.GetDirectories(os.Args[0]); err != nil {
		fmt.Println(err)
		os.Exit(8)
	}

	if conf, err = getConfiguration(rootPath, progPath); err != nil {
		fmt.Println(err)
		os.Exit(16)
	}

	if !filepath.IsAbs(conf.LogConfiguration.LogDirectory) {
		conf.LogConfiguration.LogDirectory = filepath.Join(rootPath, conf.LogConfiguration.LogDirectory)
	}
	logger.Configure("ovs_gateway", conf.LogConfiguration)

	if gate, err = ovsgate.NewInstance(conf); err != nil {
		log.Error().Err(err).Msg("Error creating gateway")
		os.Exit(16)
	}

	err = gate.Start()

	if err != nil {
		log.Error().Err(err).Msg("Error starting gateway")
		os.Exit(8)
	}

	os.Exit(0)

}

func getConfiguration(root, prog string) (config.OverseerGatewayConfig, error) {

	var err error
	var conf config.OverseerGatewayConfig

	if configFile != "" {
		if conf, err = config.Load(configFile); err != nil {
			return conf, err
		}
	} else {

		if conf, err = config.Load(filepath.Join(root, "config", "gateway.json")); err != nil {
			return conf, err
		}
	}

	if err = validator.Valid.Validate(conf); err != nil {
		return conf, err
	}

	return conf, nil
}

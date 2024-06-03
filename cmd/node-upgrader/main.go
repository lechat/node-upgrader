package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/lechat/node-upgrader/internal"
	"go.uber.org/zap"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Printf("Error parsing command-line arguments: %v\n", err)
		return 1
	}

	logLevel, logLevelSet := arguments["--log-level"].(string)
	configPath := arguments["--config-path"].(string)

	config, err := getConfig(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return 1
	}

	if !logLevelSet {
		logLevel = config.LogLevel
	}

	logger, err := internal.InitializeLogger(logLevel)
	if err != nil {
		fmt.Printf("Can't initialize zap logger: %v\n", err)
		return 1
	}
	defer logger.Sync()

	logger.Info("Starting node-upgrader",
		zap.String("log_level", logLevel),
		zap.String("config_path", configPath),
	)

	collector := internal.NewCollector(config, configPath, logger)
	collector.Start()

	return 0
}

func getConfig(configPath string) (*internal.Config, error) {
	config, err := internal.ReadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	return config, nil
}

const usage = `node-restarter

 Usage:
   node-restarter [--log-level=<level>] [--config-path=<path>]
   node-restarter -h | --help

 Options:
   -h --help              Show this screen.
   --log-level=<level>    Set the log level.
   --config-path=<path>   Set the configuration path [default: ./].
 `

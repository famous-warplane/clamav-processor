package main

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	InputPath, OutputPath, QuarantinePath string
	ScanInterval                          int
}

const (
	inputPathName      string = "inputPath"
	outputPathName     string = "outputPath"
	quarantinePathName string = "quarantinePath"
	scanIntervalName   string = "scanInterval"
)

var cfgFile string

func initConfig() (config, error) {
	setupFlags()

	err := setupViper()
	if err != nil {
		return config{}, err
	}

	return populateConfig()
}

func setupFlags() {
	pflag.StringVarP(&cfgFile, "config", "c", "", "path to a configuration file to be used (default is /etc/cap/cap.toml)")
	pflag.StringP(inputPathName, "i", "", "the directory containing incoming files to scan")
	pflag.StringP(outputPathName, "o", "", "the directory to which safe scanned files are moved")
	pflag.StringP(quarantinePathName, "q", "", "the directory to which quarantined files are moved")
	pflag.IntP("scanInterval", "s", 0, "the time interval at which the scan occurs")
	pflag.Parse()
}

func populateConfig() (config, error) {
	conf := config{}

	err := viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}

	if conf.InputPath == "" {
		return conf, errors.New("no input directory provided")
	}

	if conf.OutputPath == "" {
		return conf, errors.New("no output directory provided")
	}

	if conf.QuarantinePath == "" {
		return conf, errors.New("no quarantine directory provided")
	}

	return conf, nil
}

func setupViper() error {
	viper.BindPFlags(pflag.CommandLine)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/cap")
		viper.AddConfigPath(".")
		viper.SetConfigName("cap")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %w \n", err)
	}
	logger.Infof("Using configuration file %q", viper.ConfigFileUsed())

	viper.AutomaticEnv()

	return nil
}

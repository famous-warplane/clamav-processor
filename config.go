package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type config struct {
	inputPath, outputPath, quarantinePath string
}

func initConfig() (config, error) {
	conf := config{}
	outPath := pflag.StringP("outputPath", "o", "", "directory to which safe scanned files are moved")
	quarantinePath := pflag.StringP("quarantinePath", "q", "", "directory to which quarantined files are moved")
	pflag.Parse()
	if len(os.Args) < 2 {
		return conf, fmt.Errorf("no input directory provided. Usage: %s", UsageString)
	}
	conf.inputPath = os.Args[1]

	if *outPath == "" {
		return conf, fmt.Errorf("no output directory provided. Usage: %s", UsageString)
	}
	conf.outputPath = *outPath

	if *quarantinePath == "" {
		return conf, fmt.Errorf("no quarantine directory provided. Usage: %s", UsageString)
	}
	conf.quarantinePath = *quarantinePath
	return conf, nil

}

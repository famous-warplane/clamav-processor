package main

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type config struct {
	inputPath, outputPath, quarantinePath string
	scanInterval                          time.Duration
}

func initConfig() (config, error) {
	conf := config{}
	inPath := pflag.StringP("inputPath", "i", "", "path containing files to be scanned")
	outPath := pflag.StringP("outputPath", "o", "", "directory to which safe scanned files are moved")
	quarantinePath := pflag.StringP("quarantinePath", "q", "", "directory to which quarantined files are moved")
	interval := pflag.IntP("period", "p", 60, "time in seconds between each scan")

	pflag.Parse()

	conf.scanInterval = time.Duration(*interval) * time.Second

	if *inPath == "" {
		return conf, errors.New("no input directory provided")
	}
	conf.inputPath = os.Args[1]

	if *outPath == "" {
		return conf, errors.New("no output directory provided")
	}
	conf.outputPath = *outPath

	if *quarantinePath == "" {
		return conf, errors.New("no quarantine directory provided")
	}
	conf.quarantinePath = *quarantinePath
	return conf, nil

}

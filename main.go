package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	s "strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const UsageString string = "cap INPUT_PATH -q QUARANTINE_PATH -o OUTPUT_PATH"

var logger zap.SugaredLogger

func initLogger() {
	logConf := zap.NewProductionConfig()
	logConf.Encoding = "console"
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	unsugared, err := logConf.Build()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	logger = *unsugared.Sugar()
}

func main() {
	initLogger()
	conf, err := initConfig()
	if err != nil {
		logger.Fatalf("failed to initialise application configuation: %v", err)
	}

	scanResultsArray := runClamdscan(conf.inputPath, conf.quarantinePath)
	// Look throgh the clamscan output for files marked OK
	for _, line := range scanResultsArray {
		/*  The clamdscan outputs in a format similar to:
			'/path/to/scanned/file: scanned OK'
		so we search for lines that are OK in the output,
		then split the line on ':' to extract just the filepath
		*/
		if s.Contains(line, "OK") {
			outputParts := s.Split(line, ":")
			if len(outputParts) < 1 {
				log.Printf("Failed to parse output line %q: Line does not contain expected character ':'", outputParts)
				continue
			}
			filePath := outputParts[0]
			filePathParts := s.Split(filePath, "/")
			fileName := filePathParts[len(filePathParts)-1]

			// Move the files marked OK to a seperate 'safe' directory
			os.Rename(filePath, conf.outputPath+"/"+fileName)
		}
	}
}

func runClamdscan(scanDirectory string, quarantineDirectory string) []string {
	clamdscanCmd := "clamdscan " + scanDirectory + " --fdpass  --no-summary --move=" + quarantineDirectory

	cmd := exec.Command("/bin/sh", "-c", clamdscanCmd)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("ERROR! " + err.Error())
	}

	return s.Split(string(stdout), "\n")

}

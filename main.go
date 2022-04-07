package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger zap.SugaredLogger

func initLogger() {
	logConf := zap.NewProductionConfig()
	logConf.Encoding = "console"
	logConf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConf.DisableStacktrace = true
	unsugared, err := logConf.Build()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	logger = *unsugared.Sugar()
}

func main() {
	initLogger()
	_, err := exec.LookPath("clamdscan")
	if err != nil {
		logger.Fatal("could not find 'clamdscan' in the PATH environment variable")
	}
	conf, err := initConfig()
	if err != nil {
		fmt.Println(err)
		pflag.Usage()
		os.Exit(1)
	}

	ticker := time.NewTicker(conf.scanInterval)
	defer ticker.Stop()

	for range ticker.C {
		logger.Debug("processing input path")
		err = scan(conf)
		if err != nil {
			logger.Errorf("error while scanning directory: %v", err)
		}
	}
}

//scan runs a clamav scan and moves 'OK' files to the specified output path
func scan(conf config) error {
	logger.Debug("running scan")
	scanResultsArray, err := runClamdscan(conf.inputPath, conf.quarantinePath)
	if err != nil {
		return err
	}
	logger.Info("scan complete")
	logger.Info("processing 'OK' files")
	for _, line := range scanResultsArray {
		/*  The clamdscan outputs in a format similar to:
			'/path/to/scanned/file: scanned OK'
		so we search for lines that are OK in the output,
		then split the line on ':' to extract just the filepath
		*/
		if strings.Contains(line, "OK") {
			outputParts := strings.Split(line, ":")
			if len(outputParts) < 1 {
				logger.Warnf("Failed to parse output line %q: Line does not contain expected character ':'", outputParts)
				continue
			}
			filePath := outputParts[0]
			filePathParts := strings.Split(filePath, "/")
			fileName := filePathParts[len(filePathParts)-1]

			os.Rename(filePath, conf.outputPath+"/"+fileName)
		}
	}
	return nil
}

func runClamdscan(scanDirectory string, quarantineDirectory string) ([]string, error) {
	clamdscanCmd := "clamdscan " + scanDirectory + " --fdpass  --no-summary --move=" + quarantineDirectory

	cmd := exec.Command("/bin/sh", "-c", clamdscanCmd)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("could not attach to command output: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start clamav scan: %w", err)
	}

	cmdOut, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to process scan output: %w", err)
	}
	cmdErr, err := io.ReadAll(stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to process scan output: %w", err)
	}
	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("clamav scan failed with error: %w: %s", err, strings.TrimSpace(string(cmdErr)))
	}
	return strings.Split(string(cmdOut), "\n"), nil

}

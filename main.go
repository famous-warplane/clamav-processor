package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

const BinaryName string = "clamdscan"

func main() {
	initLogger()

	_, err := exec.LookPath(BinaryName)
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
			fileName := path.Base(filePath)
			newPath := path.Join(conf.outputPath, fileName)
			logger.Debugf("found safe file to move %q", filePath)
			os.Rename(filePath, newPath)
			logger.Infof("moved file %q to %q", filePath, newPath)
		}
	}
	return nil
}

func runClamdscan(scanDirectory string, quarantineDirectory string) ([]string, error) {
	clamdscanCmd := fmt.Sprintf("%s %s --fdpass --no-summary --move=%s", BinaryName, path.Join(scanDirectory, "*"), quarantineDirectory)
	logger.Infof("executing command %q", clamdscanCmd)
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
	cmd.Wait()
	// clamdscan returns an error code if it finds virus files so ignore the code and output to log if we get output to STDERR
	if len(cmdErr) > 0 {
		return nil, fmt.Errorf("clamav scan failed with error: %s", strings.TrimSpace(string(cmdErr)))
	}
	return strings.Split(string(cmdOut), "\n"), nil

}

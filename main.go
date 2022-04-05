package main

import (
	// "fmt"
	"os"
	"os/exec"
	s "strings"

	flag "github.com/jessevdk/go-flags"
)

func main() {
	var opts  struct {
		directoryToScan  string `short: "d" long: "dirscan" description: "Directory to be scanned by ClamAV" required: "true"`
		quarantineDirectory string `short: "q" long: "quarantinedir" desription: "Directory to send files deemed dangerous by ClamAV to be quarantined" required: "true"`
		safeDirectory string `short: "s" long: "safedir" desription: "Directory to send files deemed safe by ClamAV to be processed" required: "true"`
	}
	flag.Parse(opts)

	scanResultsArray := runClamdscan(opts.directoryToScan, opts.quarantineDirectory)
	// Look throgh the clamscan output for files marked OK
	for _, line := range scanResultsArray {
		if s.Contains(line, "OK") {
			// Get the file path for the scanned files
			filePath := s.Split(line, ":")[0]
			// Move the files marked OK to a seperate directory
			os.Rename(filePath, opts.safeDirectory)
		}
	}
}

func runClamdscan(scanDirectory string, quarantineDirectory string) []string {
	clamdscanCmd := "clamdscan " + scanDirectory + " --no-summary --fdpass --move=" + quarantineDirectory

	cmd := exec.Command("/bin/sh", "-c", clamdscanCmd)
	stdout, _ := cmd.Output()

	outputArray := s.Split(string(stdout), "\n")
	return outputArray

}

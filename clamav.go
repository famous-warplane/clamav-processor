package main

import (
	"fmt"
	"os/exec"
	s "strings"
)

func main() {
	scanResultsArray := runClamdscan("/tmp/input/*", "/tmp/virus")
	// Look throgh the clamscan output for files marked OK
	for _, line := range scanResultsArray {
		if s.Contains(line, "OK") {
			// Get the file path for the scanned files
			filePath := s.Split(line, ":")[0]
			// Move the files marked OK to a seperate directory
			mvCmd := exec.Command("/bin/sh", "-c", "mv "+filePath+" /tmp/goodfiles")
			mvCmd.Output()
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

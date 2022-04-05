package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	s "strings"
)

func main() {
	var directoryToScan string
	var quarantineDirectory string
	var safeDirectory string

	directoryToScan = os.Args[1]
	quarantineDirectory = os.Args[2]
	safeDirectory = os.Args[3]

	scanResultsArray := runClamdscan(directoryToScan, quarantineDirectory)
	// Look throgh the clamscan output for files marked OK
	for _, line := range scanResultsArray {
		/*  The clamdscan outputs in a format similar to:
			'/path/to/scanned/file: scanned OK'
		so we search for lines that are OK in the output,
		then split the line on ':' to extract just the filepath
		*/
		if s.Contains(line, "OK") {
			parts := s.Split(line, ":")
			if len(parts) < 1 {
				log.Printf("Failed to parse output line %q: Line does not contain expected character ':'", parts)
				continue
			}
			filePath := parts[0]
			fileName := s.Split(filePath, "/")[len(filePath)-1]

			// Move the files marked OK to a seperate 'safe' directory
			os.Rename(filePath, safeDirectory+"/"+fileName)
		}
	}
}

func runClamdscan(scanDirectory string, quarantineDirectory string) []string {
	clamdscanCmd := "clamdscan " + scanDirectory + " --fdpass --no-summary --move=" + quarantineDirectory

	cmd := exec.Command("/bin/sh", "-c", clamdscanCmd)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("ERROR! " + err.Error())
	}

	return s.Split(string(stdout), "\n")

}

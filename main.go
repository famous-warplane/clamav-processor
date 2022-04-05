package main

import (
	// "fmt"
	"fmt"
	"os"
	"os/exec"
	s "strings"
	// flag "github.com/jessevdk/go-flags"
)

func main() {
	// var opts struct {
	// 	directoryToScan     string `short:"d" long:"dirscan" description:"Directory to be scanned by ClamAV" required:"true"`
	// 	quarantineDirectory string `short:"q" long:"quarantinedir" desription:"Directory to send files deemed dangerous by ClamAV to be quarantined" required:"true"`
	// 	safeDirectory       string `short:"s" long:"safedir" desription:"Directory to send files deemed safe by ClamAV to be processed" required:"true"`
	// }
	// args := []string {
	// 	"-d",
	// 	"-q",
	// 	"-s",
	// }
	// flag.ParseArgs(&opts, args)
	var directoryToScan string
	var quarantineDirectory string
	var safeDirectory string

	directoryToScan = os.Args[1]
	quarantineDirectory = os.Args[2]
	safeDirectory = os.Args[3]

	scanResultsArray := runClamdscan(directoryToScan, quarantineDirectory)
	// Look throgh the clamscan output for files marked OK
	for _, line := range scanResultsArray {
		if s.Contains(line, "OK") {
			// Get the file path for the scanned files
			filePath := s.Split(line, ":")[0]
			// Move the files marked OK to a seperate directory
			os.Rename(filePath, safeDirectory)
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

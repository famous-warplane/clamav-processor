package main

import (
    "fmt"
    "os/exec"
    s "strings"
)

func main() {
    clamdscanCmd := "clamdscan /tmp/input/* --no-summary --fdpass --move=/tmp/virus"

    cmd := exec.Command("/bin/sh", "-c", clamdscanCmd)
    stdout, _ := cmd.Output()

    outputByLine := s.Split(string(stdout), "\n")

    // Print the output
    fmt.Println(string(outputByLine[0]))

    for _, line := range outputByLine {
        if s.Contains(line, "OK") {
            filePath := s.Split(line, ":")[0]
            fmt.Println(filePath)
            exec.Command("/bin/sh", "-c", "mv" + filePath + "/tmp/goodfiles")

    }

    
        
    }
}


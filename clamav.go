package main

import (
    "fmt"
    "os/exec"
    s "strings"
)

func main() {
    app  := "clamdscan /tmp/input/* --no-summary --fdpass --move=/tmp/virus"

    cmd := exec.Command("/bin/sh", "-c", app)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    outputByLine := s.Split(string(stdout), "\n")

    // Print the output
    fmt.Println(string(outputByLine[0]))

    for _, line := range outputByLine {
        if s.Contains(line, "OK"){
            fmt.Println("gotit")
    }

    
        
    }
}


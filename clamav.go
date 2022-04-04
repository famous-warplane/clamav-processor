package main

import "fmt"
import "os/exec"

func main() {
    app  := "clamdscan /tmp/input/* --no-summary --fdpass --move=/tmp/virus"

    cmd := exec.Command("/bin/sh", "-c", app)
    stdout, err := cmd.CombinedOutput()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    fmt.Println(string(stdout))
}


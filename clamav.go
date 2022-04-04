package main

import "fmt"
import "os/exec"

func main() {
    app  := "clamdscan"

    arg0 := "/tmp/input"
    arg1 := "--no-summary"
    arg2 := "--move=/tmp/virus"
    arg3 := "--fdpass"

    cmd := exec.Command(app, arg0, arg1, arg3, arg2)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    fmt.Println(string(stdout))
}


package main

import (
    "fmt"
    "os"
    "github.com/m0bsterrabb69t/bivittatus/lootutil"
    //"strings"
)

func isPrintable(err, previousErr error) bool {
    if err.Error() == previousErr.Error() {
        return false
    }
    //fmt.Println(err == previousErr)
    return true
}

func Run() {
    cmdArgs := os.Args[1:]

    url := &cmdArgs[len(cmdArgs)-1]
    previousErr := fmt.Errorf("")
    for {
        err := lootutil.Loot(url)
        if err != nil {
            if isPrintable(err, previousErr) {
                fmt.Printf("\n\n%v\n\n", err)
                previousErr = err
            }
        } else {
            break
        }
    }

    fmt.Printf("Done.\a\n")
}

func main() {
    Run()
}

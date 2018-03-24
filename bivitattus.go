package main

import (
    "fmt"
    "os"
    "github.com/m0bsterrabb69t/bivitattus/lootutil"
    //"strings"
)

func errorExists(Errors []string, err error) bool {
    for _, v := range Errors {
        if v == err.Error() {
            return true
        }
    }
    Errors = append(Errors, err.Error())
    return false
}

func Run() {
    cmdArgs := os.Args[1:]

    url := &cmdArgs[len(cmdArgs)-1]
    Errors := make([]string, 0)
    for {
        err := lootutil.Loot(url)
        if err != nil {
            if !errorExists(Errors, err) {
                fmt.Printf("%v\n\n", err)
            }
        } else {
            break
        }
    }
}

func main() {
    Run()
}

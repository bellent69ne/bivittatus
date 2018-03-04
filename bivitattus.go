package main

import (
    "fmt"
    "os"
    "github.com/m0bsterrabb69t/bivitattus/lootutil"
)

func Run() {
    cmdArgs := os.Args[1:]

    err := lootutil.Loot(&cmdArgs[len(cmdArgs)-1])
    if err != nil {
        fmt.Printf("Error... %v\n", err)
    }
}

func main() {
    Run()
}

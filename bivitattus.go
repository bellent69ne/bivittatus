package main

import (
    "os"
    "github.com/m0bsterrabb69t/bivitattus/lootutil"
)

func Run() {
    cmdArgs := os.Args[1:]

    lootutil.Loot(&cmdArgs[len(cmdArgs)-1])
}

func main() {
    Run()
}

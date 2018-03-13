package main

import (
    "fmt"
    "os"
    "github.com/m0bsterrabb69t/bivitattus/lootutil"
    //"strings"
)

//func fileName(url *string) string {
//    result := strings.Split(*url, "/")
//
//    return result[len(result) - 1]
//}

func Run() {
    cmdArgs := os.Args[1:]

    url := &cmdArgs[len(cmdArgs)-1]
    err := lootutil.Loot(url)
    if err != nil {
        fmt.Printf("Error... %v\n", err)
    }
    //fileName(url)
}

func main() {
    Run()
}

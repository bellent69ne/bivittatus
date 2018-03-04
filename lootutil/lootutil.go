// package lootutil contains utility functions
// for web loot operations
package lootutil

import (
    "net/http"
    "io/ioutil"
    "strconv"
    "fmt"
    "errors"
)

/***************YOU SHOULD RECONSIDER THIS CODE******************/
// inspect determines support for partial requests and size of file
func inspect(url *string) (size int, err error) {
    // make HEAD request
    resp, err := http.Head(*url)
    // Handle errors if any
    if err != nil {
        return
    }
    // Check support for partial requests
    _, ok := resp.Header["Accept-Ranges"]
    if !ok {
        err = errors.New(*url + " doesn't support partial requests...")
        return
    }
    fileSize, _ := resp.Header["Content-Length"]
    // Convert file size from string to int
    size, err = strconv.Atoi(fileSize[0])

    return
}
// Create http request
func createRequest(method, block string, url *string) (
    req *http.Request, err error) {
    req, err = http.NewRequest("GET", *url, nil)
    // if error happened, return immediately
    if err != nil {
        return
    }
    // Add "Range" header to http request
    req.Header.Add("Range", block)

    return
}

const (
    KB = 1024
    CHUNK = 256 * KB
)

// Try returning response body
// lootChunk downloads piece of file
func lootChunk(url *string, startAt int) error {
    block := fmt.Sprintf("bytes=%v-%v", startAt, startAt + CHUNK)

    req, err := createRequest("GET", block, url)
    // Handle errors if any
    if err != nil {
        return err
    }

    // Create http client
    client := &http.Client{}
    // Make request
    resp, err := client.Do(req)

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    fmt.Println(body)

    return nil
}

func Loot(url *string) error {
    size, err := inspect(url)
    fmt.Println("Size = ", size)

    if err != nil {
        return err
    }

    err = lootChunk(url, 0)
    if err != nil {
        fmt.Printf("Error... %v\n", err)
    }

    return nil
}

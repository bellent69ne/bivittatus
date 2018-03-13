// package lootutil contains utility functions
// for web loot operations
package lootutil

import (
    "net/http"
    "io/ioutil"
    "strconv"
    "fmt"
    "os"
    "strings"
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
    CHUNK = 128 * KB
)

type DataChunk struct {
    data []uint8
    err error
}

// Try returning response body
// lootChunk downloads piece of file
func lootChunk(url *string, startAt int64, stream chan DataChunk) {
    block := fmt.Sprintf("bytes=%v-%v", startAt, startAt + CHUNK)

    req, err := createRequest("GET", block, url)
    // Handle errors if any
    if err != nil {
        stream <- DataChunk{nil, err}
    }

    // Create http client
    client := &http.Client{}
    // Make request
    resp, err := client.Do(req)

    if err != nil {
        stream <- DataChunk{nil, err}
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    stream <- DataChunk{body, err}
    fmt.Println("I think I got something")
}

func filename(url *string) string {
    result := strings.Split(*url, "/")

    return result[len(result) - 1]
}

func writeToFile(filename string, stream chan DataChunk) (err error) {
    received := <-stream
    if received.err != nil {
        return received.err
    }

    var file *os.File
    if _, err = os.Stat(filename); os.IsNotExist(err) {
        file, err = os.Create(filename)
        fmt.Println("Creating new file")
    } else {
        file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
        fmt.Println("Writing to the same file")
    }
    defer file.Close()

    if err != nil {
        return
    }

    file.Write(received.data)
    file.Sync()
    return
}

// Loot downloads pieces of file from url
// and assembles them into one file simultaneously
func Loot(url *string) error {
    size, err := inspect(url)

    if err != nil {

        return err
    }

    stream := make(chan DataChunk)

    var downloaded int64

    fi, err := os.Stat(filename(url))
    if err == nil {
        downloaded = fi.Size()
    } else if os.IsNotExist(err) {
        downloaded = 0
    } else if err != nil && !os.IsNotExist(err) {
        return err
    }

    for downloaded != int64(size) {
        go lootChunk(url, downloaded, stream)
        err = writeToFile(filename(url), stream)

        if err != nil {
            return err
        }

        downloaded += CHUNK
        // ALL THIS CODE STILL REQUIERS MODIFICATION
        if downloaded > int64(size) {
            downloaded = in64(size)
        }
    }

    //err = lootChunk(url, 0)
    if err != nil {
        return err
    }

    return nil
}

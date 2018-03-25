// package lootutil contains utility functions
// for web loot operations
package lootutil

import (
    "net/http"
    "io/ioutil"
    "fmt"
    "time"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "github.com/pkg/errors"
)

const (
    KB = 1024
    MB = KB * KB
    GB = MB * MB
    CHUNK = 16 * KB
)
/***************YOU SHOULD RECONSIDER THIS CODE******************/
// inspect determines support for partial requests and size of file
func inspect(url *string) (size int, err error) {
    // make HEAD request
    resp, err := http.Head(*url)
    // Handle errors if any
    if err != nil {
        return 0, errors.Wrap(err, "Failed making http head request")
    }
    // Check support for partial requests
    _, ok := resp.Header["Accept-Ranges"]
    if !ok {
        err = fmt.Errorf("%s doesn't support partial requests...", *url)
        return 0, err
    }
    fileSize, _ := resp.Header["Content-Length"]
    // Convert file size from string to int
    size, err = strconv.Atoi(fileSize[0])
    // Probably need to wrap err here////////////////////
    return
}
// Create http request
func createRequest(method, block string, url *string) (
    req *http.Request, err error) {
    req, err = http.NewRequest("GET", *url, nil)
    // if error happened, return immediately
    if err != nil {
        return nil, errors.Wrap(err, "Failed creating http request")
    }
    // Add "Range" header to http request
    req.Header.Add("Range", block)

    return
}


type DataChunk struct {
    data []uint8
    time.Duration
    err error
}

// Try returning response body
// lootChunk downloads piece of file
func lootChunk(url *string, startAt int64, stream chan DataChunk) {
    block := fmt.Sprintf("bytes=%v-%v", startAt, startAt + (CHUNK - 1))

    req, err := createRequest("GET", block, url)
    // Handle errors if any
    if err != nil {
        stream <- DataChunk{nil, 0, err}
        return
    }

    // Create http client
    client := &http.Client{}
    // start measuring time
    start := time.Now()
    // Make request
    resp, err := client.Do(req)
    // if error happened, send nil data and error to the channel
    if err != nil {
        stream <- DataChunk{nil, 0, err}
        return
    }
    defer resp.Body.Close()
    // Read all bytes from the body of response
    body, err := ioutil.ReadAll(resp.Body)
    elapsed :=time.Now().Sub(start)
    // Send that body to a channel
    stream <- DataChunk{body, elapsed, err}
}
// filename is getting a filename
func filename(url *string) string {
    result := strings.Split(*url, "/")

    return result[len(result) - 1]
}

func accessFile(filename *string) (file *os.File, err error) {
    // Create file object to open an existing file or create new one
    if _, err = os.Stat(*filename); os.IsNotExist(err) {
        file, err = os.Create(*filename)
    } else {
        // Open an existing file and move to the end of the file
        file, err = os.OpenFile(*filename, os.O_WRONLY|os.O_APPEND, 0644)
    }
    return
}

// writeToFile writes all data from a channel to a file
func writeToFile(filename string, stream chan DataChunk) (
    time.Duration, error) {
    // receive data from stream, handle any errors
    received := <-stream
    if received.err != nil {
        return 0, errors.Wrap(received.err, "Unable to get data from web")
    }

    file, err := accessFile(&filename)
    if err != nil {
        return 0, errors.Wrap(err, "Couldn't access the file")
    }
    defer file.Close()
    // Write all data to the file
    _, err = file.Write(received.data)
    if err != nil {
        return 0, errors.Wrap(err, "Failed writing data to a file")
    }
    file.Sync()
    return received.Duration, nil
}

func ttyWidth() (width int, err error) {
    cmd := exec.Command("stty", "size")
    cmd.Stdin = os.Stdin
    out, err := cmd.Output()

    if err != nil {
        return 0, errors.Wrap(err, "Failed calculating tty width")
    }

    out = out[:len(out)-1]
    strOut := string(out)
    splitted := strings.Split(strOut, " ")

    width, err = strconv.Atoi(splitted[1])
    return
}

func state(percent int) string {
    width, err := ttyWidth()
    if err != nil {
        fmt.Println(err)
        return ""
    }
    totalLength := width * 35 / 100
    stateLength := totalLength * percent / 100

    state := make([]rune, totalLength + 2)
    state[0] = '|'
    for i := 1; i <= totalLength; i++ {
        if i == stateLength {
            state[i] = '>'
        } else if i < stateLength {
            state[i] = '='
        } else {
            state[i] = ' '
        }
    }
    state[len(state) - 1] = '|'
    return string(state)
}

func printStatus(nextChunk int64, size int, elapsed time.Duration) {
    fileSize := 0
    var strSize, strGot string
    switch {
    case size > KB && size < MB:
        {
            fileSize = size / KB
            strSize = fmt.Sprintf("%dkB", fileSize)
            strGot = fmt.Sprintf("%.2fkB", float64(nextChunk) / KB)
        }

    case size > MB && size < GB:
        {
            fileSize = size / MB
            strSize = fmt.Sprintf("%dmB", fileSize)
            strGot = fmt.Sprintf("%.2fmB", float64(nextChunk) / MB)
        }
    }
    speed := CHUNK / (elapsed / time.Millisecond)
    speed *= 1000
    speed /= 1024
    percent := int(float64(nextChunk) / float64(size) * 100)
    //fmt.Printf("\r    %s    %s    %dkB/s    %s %d%% ",
     //    strSize, strGot, int(speed), state(percent), percent)
     fmt.Printf("\r    %s %d%%    %s    %dkB/s    %s   ",
        state(percent), percent, strGot, int(speed), strSize)
}

// Loot downloads pieces of file from url
// and assembles them into one file simultaneously
func Loot(url *string) error {
    size, err := inspect(url)

    if err != nil {
        return err
    }

    stream := make(chan DataChunk)

    var nextChunk int64

    fi, err := os.Stat(filename(url))
    switch {
    case err == nil:
        nextChunk = fi.Size()
    case os.IsNotExist(err):
        nextChunk = 0
    case err != nil && !os.IsNotExist(err):
        return err
    }

    fmt.Printf("\nFile: %s\n", filename(url))
    for nextChunk != int64(size) {
        go lootChunk(url, nextChunk, stream)
        elapsed, err := writeToFile(filename(url), stream)

        if err != nil {
            return err
        }

        nextChunk += CHUNK
        // ALL THIS CODE STILL REQUIERS MODIFICATION
        if nextChunk > int64(size) {
            nextChunk = int64(size)
        }
        //fmt.Printf("\telapsed %v\t%v\n", elapsed, int(elapsed))
        printStatus(nextChunk, size, elapsed)
    }
    fmt.Printf("\n\n")


    return err
}

package lootutil

import (
    "os"
    "time"
    "github.com/pkg/errors"
    "strings"
)
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

// package lootutil contains utility functions
// for web loot operations
package lootutil

import (
    "net/http"
    "io/ioutil"
    "fmt"
    "errors"
)

// Determines support for partial requests
func canBeDivided(url *string) bool {
    // make HEAD request
    resp, err := http.Head(*url)
    // Handle errors if any
    if err != nil {
        fmt.Printf("Error... %v\n", err)
        return false
    }
    // Check support for partial requests
    _, ok := resp.Header["Accept-Ranges"]

    return ok
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

    return nil
}

func Loot(url *string) error {
    ok := canBeDivided(url)

    if !ok {
        return errors.New("Cannot download the file!!!")
    }

    err := lootChunk(url, 0, 100000)
    if err != nil {
        fmt.Printf("Error... %v\n", err)
    }

    /*resp, err := http.Get(*url)
    if err != nil {
        fmt.Printf("Error occured: %v\n", err)
        return err
    }*/ /*
    req, err := http.NewRequest("GET", *url, nil)
    const chunk = 100 * 1024
    req.Header.Add("Range", "bytes=0-100")// + string(chunk))

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        fmt.Printf("Error occured: %v\n", err)
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        fmt.Printf("Error occured: %v\n", err)
        return err
    }
*/
    //fmt.Println("Message body length: ", len(body))

    //fmt.Println("Content-Type: ", resp.Header["Content-Type"])
    //fmt.Println("Range: ", req.Header["Range"])
    //fmt.Println("Accept-Ranges: ", resp.Header["Accept-Ranges"])
    //fmt.Println("Content-Length: ", resp.Header["Content-Length"])
    return nil
}

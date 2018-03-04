// package lootutil contains utility functions
// for web loot operations
package lootutil

import (
    "net/http"
    //"io/ioutil"
    "fmt"
)

func canBeDivided(url *string) (ok bool, err error) {
    resp, err := http.Head(*url)
    _, ok = resp.Header["Accept-Ranges"]

    return
}

func lootChunk(url *string, start, end int) error {
    req, err := http.NewRequest("GET"m *url, nil)
    if err != nil {
        return err
    }

    range := fmt.Sprintf("bytes=%v-%v", start, end)
    req.Header.Add("Range", )
}

func Loot(url *string) error {
    ok, err := canBeDivided(url)
    if err != nil {
        return err
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

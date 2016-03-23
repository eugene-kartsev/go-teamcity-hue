package hue

import (
    "net/http"
    "bytes"
    "time"
    "github.com/eugene-kartsev/go-teamcity-hue/server/httpHeaders"
    "io/ioutil"
    "errors"
)

func Signal(url string, signal int) (ok bool, err error) {
    if signal == RED {
        return setRed(url)
    }
    if signal == GREEN {
        return setGreen(url)
    }

    return false, errors.New("Signal is unknown")
}

func setColor(url string, requestBody []byte) (ok bool, err error) {
    ok = false

    request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))

    if err != nil {
        return
    }

    request.Header.Set(httpHeaders.CONTENT_TYPE, httpHeaders.APPLICATION_JSON)

    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := client.Do(request)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    _, err = ioutil.ReadAll(resp.Body)

    ok = err == nil
    return
}


func setGreen(url string) (bool, error) {
    result := make(chan error)

    go func() {
        bytes := []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`)
        _, err := setColor(url + "1/state", bytes)
        result <- err
    }()
    go func() {
        bytes := []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`)
        _, err := setColor(url + "2/state", bytes)
        result <- err
    }()

    err1 := <-result
    err2 := <-result

    if err1 != nil {
        return false, err1
    }

    if err2 != nil {
        return false, err2
    }

    return true, nil
}

func setRed(url string) (bool, error) {
    result := make(chan error)

    go func() {
        _, err := setColor(url + "1/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
        result <- err
    }()
    go func() {
        _, err := setColor(url + "2/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
        result <- err
    }()

    err1 := <- result
    err2 := <- result

    if err1 != nil {
        return false, err1
    }

    if err2 != nil {
        return false, err2
    }

    return true, nil
}

package hue

import (
    "net/http"
    "bytes"
    "time"
    "github.com/eugene-kartsev/go-teamcity-hue/server/httpHeaders"
    "io/ioutil"
    "errors"
)

var greenBytes = []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`)
var redBytes   = []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`)

func Signal(url string, signal int) (bool, error) {
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
    return setColor(url, greenBytes)
}

func setRed(url string) (bool, error) {
    return setColor(url, redBytes)
}

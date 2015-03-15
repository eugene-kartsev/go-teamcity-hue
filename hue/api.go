package hue

import (
    "net/http"
    "io/ioutil"
    "fmt"
    "bytes"
)


var _apiUrl string
func Init(apiUrl string) {
    _apiUrl = apiUrl

    fmt.Println(_apiUrl)
}

func setColor(url string, requestBody []byte) {

    fmt.Println(url)

    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("ERROR: " + err.Error())
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func Green() {
    fmt.Println("GREEN")
    go setColor(_apiUrl + "1/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`))
    go setColor(_apiUrl + "2/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`))
}

func Red() {
    fmt.Println("RED")
    go setColor(_apiUrl + "1/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
    go setColor(_apiUrl + "2/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
}

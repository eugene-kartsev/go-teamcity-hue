package hue

import (
    "go-teamcity-hue/config"
    "fmt"
    "net/http"
    "bytes"
    "strconv"
    "time"
)

type Api interface {
    Signal(signalType int)
}

const (
    RED   int = 0
    GREEN int = 1
)

type api struct {
    config config.IConfig
}

func Create(cfg config.IConfig) (Api, error) {
    _api := api {
        config: cfg,
    }

    return &_api, nil
}

func (self *api) Signal(signal int) {
    fmt.Println("hue. Signal: " + strconv.Itoa(signal))

    if signal == RED {
        self.setRed()
    }
    if signal == GREEN {
        self.setGreen()
    }
}

func setColor(url string, requestBody []byte) {
    fmt.Println(url)

    fmt.Println("hue. Request url: " + url)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))

    if err != nil {
        fmt.Println("ERROR: " + err.Error())
        return
    }
    defer req.Body.Close()

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("ERROR: " + err.Error())
        return
    }
    defer resp.Body.Close()

    //body, err := ioutil.ReadAll(resp.Body)
}


func (self *api) setGreen() {
    go setColor(self.config.GetHueApiUrl() + "1/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`))
    go setColor(self.config.GetHueApiUrl() + "2/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":22784}`))
}

func (self *api) setRed() {
    go setColor(self.config.GetHueApiUrl() + "1/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
    go setColor(self.config.GetHueApiUrl() + "2/state", []byte(`{"on":true, "sat":249, "bri":113, "hue":1000}`))
}

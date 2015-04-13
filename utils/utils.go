package utils
import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
)

const (
    REQUEST_TIMEOUT = 5 * time.Second
)

func ReadJsonWithCredentials(url string, obj interface{}, login string, password string) error {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    req.SetBasicAuth(login, password)
    req.Header.Add("Accept", "application/json")

    return readJsonInternal(req, obj)
}

func ReadJson(url string, obj interface{}) error {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    return readJsonInternal(req, obj)
}

func readJsonInternal(req *http.Request, obj interface{}) error {
    client := &http.Client{
        Timeout : REQUEST_TIMEOUT,
    }

    res, err := client.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    bytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }

    err = json.Unmarshal(bytes, obj)
    if err != nil {
        return err
    }

    return nil
}
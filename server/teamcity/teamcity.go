package teamcity

import (
    "crypto/tls"
    "crypto/rand"
    "net/http"
    "net"
    "time"
    "strings"
    "compress/gzip"
    "io/ioutil"
    "encoding/json"
    "errors"
    "github.com/eugene-kartsev/go-teamcity-hue/server/httpHeaders"
)

func GetBuildStatus(url string, login string, password string) (buildStatus string, err error) {
    buildStatus = ERROR
    ssl := &tls.Config{
        InsecureSkipVerify: true,
    }
    ssl.Rand = rand.Reader

    client := &http.Client{
        Transport: &http.Transport{
            Dial: func(network, addr string)(net.Conn, error) {
                return net.DialTimeout(network, addr, time.Duration(time.Second * 3))
            },
            TLSClientConfig: ssl,
        },
    }

    requestBody := strings.NewReader("{}");
    req, err := http.NewRequest("GET", url, requestBody)

    if err != nil {
        return
    }

    // some custom headers
    req.SetBasicAuth(login, password)
    req.Header.Add(httpHeaders.CONTENT_TYPE, httpHeaders.APPLICATION_JSON)
    req.Header.Add(httpHeaders.ACCEPT, httpHeaders.APPLICATION_JSON)

    response, err := client.Do(req)
    if err != nil {
        return
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        err = errors.New("FAILED: response.Status: " + response.Status)
        return
    }

    // in case response is gziped, run that through the gzip reader to de decompress
    if strings.Contains(response.Header.Get(httpHeaders.CONTENT_ENCODING), httpHeaders.GZIP) {
        response.Body, err = gzip.NewReader(response.Body)
        if err != nil {
            return
        }
    }

    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }

    var obj buildList

    str := string(data)
    bytes := []byte(str)
    err = json.Unmarshal(bytes, &obj)

    if err != nil {
        return
    }

    if obj.Builds == nil || len(obj.Builds) < 1 {
        err = errors.New("Could not parse string: " + str)
        return
    }

    buildStatus = obj.Builds[0].Status

    return buildStatus, nil
}

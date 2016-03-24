package config

import (
    "io/ioutil"
    "os"
    "errors"
    "encoding/json"
)

func Read() (*Config, error) {
    var cfg Config

    if !exists() {
        err := create()
        if err != nil {
            return nil, errors.New("Cannot create a configuration file. " + err.Error())
        }

        return nil, errors.New("Configuration file has been created.\n\rModify configuration file and restart the application.")
    }

    bytes, err := ioutil.ReadFile(configFileName)
    if err != nil {
        return nil, errors.New("Cannot read a configuration file. " + err.Error())
    }


    err = json.Unmarshal(bytes, &cfg)

    if err != nil {
        return nil, err
    }

    return &cfg, nil
}

func exists() bool {
    _, err := os.Stat(configFileName);
    return err == nil || os.IsExist(err);
}

func create() error {
    file, err := os.OpenFile(configFileName, os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {
        return err
    }
    defer file.Close()

    file.WriteString(configFileTemplate)
    return nil
}

const (
    configFileName = "config"
    configFileTemplate =
    `
    {
        "version": "0.01",
        "hueNodes": [{
            "id": "hue1",
            "url": "<HUE_URL>"
        }],
        "teamcityNodes": [{
            "id": "tc1",
            "url": "<TEAMCITY_BUILD_URL>",
            "login": "<USER_LOGIN>",
            "password": "<USER_PASSWORD>",
            "interval": 10
        }],
        "map": [{
            "hueId": "hue1",
            "teamcityIds": ["tc1"]
        }]
    }
    `
)

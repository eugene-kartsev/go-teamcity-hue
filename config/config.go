package config

import (
    "flag"
    "errors"
)

type IConfig interface {
    GetTeamcityApiUrl() string
    GetTeamcityCredentials() (string, string)
    GetHueApiUrl() string
}

var _config mainConfig = mainConfig{}
func init() {
    _config.teamcityApiUrl   = flag.String("tc.url", "", "TeamCity api url")
    _config.teamcityLogin    = flag.String("tc.login", "", "TeamCity api Login")
    _config.teamcityPassword = flag.String("tc.password", "", "TeamCity api Password")

    _config.hueApiUrl = flag.String("hue.url", "", "HUE api url")
}

func Create() (IConfig, error) {
    flag.Parse()

    if *_config.teamcityApiUrl == "" || *_config.hueApiUrl == "" {
        err := errors.New("Error: config has not been initialized.")
        flag.PrintDefaults()

        return nil, err
    }

    return &_config, nil
}


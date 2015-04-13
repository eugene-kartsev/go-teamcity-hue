package config


type mainConfig struct {
    teamcityApiUrl *string
    teamcityLogin *string
    teamcityPassword *string

    hueApiUrl *string
}

func (self mainConfig) GetTeamcityApiUrl() string {
    return *self.teamcityApiUrl
}

func (self mainConfig) GetTeamcityCredentials() (string, string) {
    return *self.teamcityLogin, *self.teamcityPassword
}

func (self mainConfig) GetHueApiUrl() string {
    return *self.hueApiUrl
}

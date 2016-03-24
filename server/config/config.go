package config

type Config struct {
    Version       string
    HueNodes      []HueNode      `json:"hueNodes"`
    TeamCityNodes []TeamCityNode `json:"teamCityNodes"`
    Map           []Mapping
}

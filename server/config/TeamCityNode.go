package config

type TeamCityNode struct {
    Id       string
    Url      string
    Login    string
    Password string
    HueNodes []string
    Interval int
}
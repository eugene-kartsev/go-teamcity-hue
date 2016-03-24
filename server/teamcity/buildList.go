package teamcity

type buildList struct {
    Count int
    Href  string
    Builds []build `json:"build"`
}

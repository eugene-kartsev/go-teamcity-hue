package teamcity

type buildList struct {
    Count int
    Href  string
    Builds []build `json:"build"`
}

type build struct {
    Id          int
    BuildTypeId string
    Number      string
    Status      string
    State       string
}

const (
    SUCCESS = "SUCCESS"
    FAILURE = "FAILURE"
    ERROR   = "ERROR"
)

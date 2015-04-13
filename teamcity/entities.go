package teamcity

type Project struct {
    Id              string
    Name            string
    Description     string
    Href            string
    WebUrl          string
    ParentProjectId string
    BuildTypes      []BuildType `json:"-"`
}

type BuildType struct {
    Id          string
    Name        string
    ProjectName string
    ProjectId   string
    Href        string
    WebUrl      string
}

type projectList struct {
    Count   int
    Projects []Project `json:"project"`
}

type projectBuildTypes struct {
    Id         string
    BuildTypes buildTypeList
}

type buildTypeList struct {
    Count      int
    BuildTypes []BuildType `json:"BuildType"`
}

type buildList struct {
    Count int
    Href  string
    Builds []Build `json:"Build"`
}

type Build struct {
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
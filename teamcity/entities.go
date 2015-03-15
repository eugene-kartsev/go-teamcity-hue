package teamcity

import "time"


//items to track:
//project/build_type
//examples:
//  */* - every single project and build type
//  PressReader/* - everything from Pressreader project
//  */pressreader-common - find all build types with name 'pressreader-common'
//  PressReader/pressreader-common - only a single project
//  !*/* - this doesn't make sense, skip this rule
//  !PressReader/pressreader-common - exclude
//  !PressReader/* - exclude
//  !*/pressreader-common - exclude

type Config struct {
    Login         string
    Password      string
    RefreshInSec  time.Duration
}

type Project struct {
    Id              string
    Name            string
    Description     string
    Href            string
    WebUrl          string
    ParentProjectId string
    BuildTypes      BuildTypeList
}

type ProjectList struct {
    Count   int
    Project []Project
}

type BuildType struct {
    Id          string
    Name        string
    ProjectName string
    ProjectId   string
    Href        string
    WebUrl      string
}

type BuildTypeList struct {
    Count     int
    BuildType []BuildType
}

const (
    SUCCESS = "SUCCESS"
    FAILURE = "FAILURE"
    ERROR   = "ERROR"
)

type Build struct {
    Id          int
    BuildTypeId string
    Number      string
    Status      string
    State       string
}

type BuildList struct {
    Count int
    Href  string
    Build []Build
}

type buildState struct {
    buildTypeId  string
    currentState bool
    oldState     bool
    isNew        bool
}

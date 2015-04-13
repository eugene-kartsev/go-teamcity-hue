package dispatcher
import (
    "go-teamcity-hue/teamcity"
    "go-teamcity-hue/hue"
    "time"
    "fmt"
    "strconv"
)

const (
    REFRESH_PROJECTS_INTERVAL_SEC = 60 * time.Second
    REFRESH_STATUS_INTERVAL_SEC = 15 * time.Second
)

const (
    STATUS_UNDEFINED = 0
    STATUS_SUCCESS = 1
    STATUS_FAILED = 2
)

var canReed      chan bool
var tc           teamcity.Api
var hueLights    hue.Api

var _projects    map[string] teamcity.Project
var _buildState  map[string] buildState

func init() {
    canReed = make(chan bool, 1)
}

func Start(teamcityApi teamcity.Api, hue hue.Api) {
    tc = teamcityApi
    hueLights = hue

    canReed <- true
    go func() {
        for {
            time.Sleep(REFRESH_PROJECTS_INTERVAL_SEC)
            fetchProjectStructure()
        }
    }()

    go func() {
        fetchProjectStructure() //first time only
        for {
            newGlobalStatus := fetchBuildStatus()
            updateBuildStatus(newGlobalStatus)
            time.Sleep(REFRESH_STATUS_INTERVAL_SEC)
        }
    }()
}

func projectMapHasChanges(projects map[string] teamcity.Project) bool {
    hasChanges := false

    <- canReed

    hasChanges = len(_projects) != len(projects)

    if !hasChanges {
        for projectId := range _projects {
            _p := _projects[projectId]
            p, hasValue := projects[projectId]

            hasChanges = hasChanges ||
                !hasValue ||
                _p.Id              != p.Id ||
                _p.ParentProjectId != p.ParentProjectId ||
                _p.WebUrl          != p.WebUrl

            if !hasChanges && len(p.BuildTypes) == len(_p.BuildTypes) && len(p.BuildTypes) > 0 {
                for _, _buildType := range p.BuildTypes {
                    contains := false
                    for _, buildType := range _p.BuildTypes {
                        contains = contains || _buildType.Id == buildType.Id
                    }
                    hasChanges = hasChanges || !contains
                }
            }
        }
    }

    canReed <- true

    return hasChanges
}


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
func matchBuild(build buildState) bool {
    return build.projectName == "PressReader"
}

func fetchBuildStatus() *[]buildState {
    <- canReed
    buildsToFetch := []buildState{}

    for _, build := range _buildState {
        if matchBuild(build) {
            buildsToFetch = append(buildsToFetch, build)
        }
    }

    canReed <- true

    done := make(chan bool, len(buildsToFetch))
    for index, _ := range buildsToFetch {
        build := &buildsToFetch[index]
        go func() {
            state, err := tc.GetBuildState(&build.buildId)
            if err != nil && state != nil {
                success := state.State == teamcity.SUCCESS
                if success {
                    build.status = STATUS_SUCCESS
                }
            }
            done <- true
        }()
    }
    for range buildsToFetch {
        <- done
    }
    return &buildsToFetch
}

func updateBuildStatus(newState *[]buildState) {
    <- canReed

    newStatus := STATUS_SUCCESS

    for _, newBuildState := range (*newState) {
        oldBuildState, ok := _buildState[newBuildState.buildId]
        if ok && newBuildState.status != oldBuildState.status {
            if newStatus != STATUS_FAILED {
                newStatus = newBuildState.status
            }
        }
    }

    canReed <- true

    if newStatus == STATUS_SUCCESS {
        go hueLights.Signal(hue.GREEN)
    } else {
        go hueLights.Signal(hue.RED)
    }
    fmt.Println("TeamCity. New Status: " + strconv.Itoa(newStatus))
}

func fetchProjectStructure() {
    projects, err := loadProjects()

    if err == nil {
        if projectMapHasChanges(projects) {
            <- canReed

            _projects    = projects
            _buildState = make(map[string] buildState)

            for _, project := range _projects {
                for _, buildType := range project.BuildTypes {
                    _buildState[buildType.Id] = buildState{
                        projectId:   project.Id,
                        projectName: project.Name,
                        buildId:     buildType.Id,
                        buildName:   buildType.Name,
                        status:      STATUS_UNDEFINED,
                        oldStatus:   STATUS_UNDEFINED,
                    }
                }
            }

            canReed <- true
        }
    }
}

func loadProjects() (map[string] teamcity.Project, error) {
    loadBuildTypes := true
    projects, err := tc.GetProjects(loadBuildTypes)

    if err != nil {
        return nil, err
    }

    result := make(map[string] teamcity.Project)

    for _, project := range projects {
        result[project.Id] = project
    }

    return result, nil
}

type buildState struct {
    projectId   string
    projectName string
    buildId     string
    buildName   string
    status      int
    oldStatus   int
}
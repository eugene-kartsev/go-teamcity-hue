package teamcity

import (
    "fmt"
    "net/http"
    "time"
    "io/ioutil"
    "encoding/json"
)

type Api interface {
    Watch(onStateChanged func(statusOK bool))
}

type api struct {
    config        Config
    apiUrl        string
    globalStateChanged chan bool
    buildStateChanged  chan buildState
    buildStates        map[string] *buildState
    canReed            chan bool
}

func Setup(config Config) (Api) {
    self := new(api)
    self.config = config
    self.apiUrl = "http://ndwdbld001.se.newspaperdirect.com/httpAuth/app/rest/"

    self.buildStates = make(map[string] *buildState)

    self.globalStateChanged = make(chan bool)
    self.buildStateChanged  = make(chan buildState)
    self.canReed            = make(chan bool, 1)

    self.canReed <- true

    return self
}

func (self api) Watch(onStateChanged func(statusOK bool)) {
    go self.stateWatcher(onStateChanged)
    go self.stateUpdater()

    go self.newWorker()
}

func (self api) stateWatcher(onStateChanged func(ok bool)) {
    for {
        newState := <-self.globalStateChanged
        fmt.Println("newState:")
        fmt.Println(newState)
        go onStateChanged(newState)
    }
}

func (self api) stateUpdater() {
    for {
        newState := <- self.buildStateChanged
        fmt.Println("stateUpdater: START")
        <- self.canReed

        if item := self.buildStates[newState.buildTypeId]; item == nil {
            newState.oldState = newState.currentState;
            newState.isNew = true
            self.buildStates[newState.buildTypeId] = &newState;
        } else {
            item.currentState = newState.currentState
        }
        fmt.Println("stateUpdater: DONE")

        self.canReed <- true
    }
}

func (self api) checkState() {
    changed     := false
    globalState := true

    <- self.canReed
    fmt.Println("checkState: START")
    for _, buildState := range self.buildStates {
        fmt.Print(buildState.currentState)
        fmt.Print(" - ")
        fmt.Print(buildState.oldState)
        fmt.Println()
        if buildState.currentState != buildState.oldState {
            changed = true
            buildState.oldState = buildState.currentState
        }
        if(buildState.isNew) {
            changed = true
            buildState.isNew = false
        }
        if globalState && !buildState.currentState  {
            globalState = false
        }
    }

    fmt.Println("checkState: DONE")
    fmt.Println(changed)
    fmt.Println(globalState)

    self.canReed <- true

    if changed {
        self.globalStateChanged <- globalState
    }
}

func (self api) newWorker() {
    for {
        self.worker()
        time.Sleep(self.config.RefreshInSec)
    }
}

func (self api) worker() {
    projects := self.getProjects()

    done := make(chan bool, len(projects))
    for _, p := range projects {
        go (func(project Project) {
            self.processProject(project)
            done<-true
        })(p);
    }
    for range projects {
        <-done
    }
    fmt.Println("DONE ALL PROJECTS")
    go self.checkState()
}

func (self api) getProjects() []Project {
    url := self.apiUrl + "projects"

    fmt.Println(url)
    var content ProjectList
    if err := self.readJson(url, &content); err != nil {
        fmt.Println("ERROR " + err.Error())
        return nil
    }

    return content.Project
}

func (self api) getBuildTypes(project Project) []BuildType {
    url := self.apiUrl + "projects/id:" + project.Id

    var content Project
    if err := self.readJson(url, &content); err != nil {
        fmt.Println("ERROR: " + err.Error())
        return nil
    }

    return content.BuildTypes.BuildType
}

func (self api) getBuilds(buildType BuildType) []Build {
    url := self.apiUrl + "buildTypes/id:"+ buildType.Id +"/builds?count=1&start=0&running=false&canceled=false"
    var content BuildList
    if err := self.readJson(url, &content); err != nil {
        fmt.Println("ERROR: " + err.Error())
        return nil
    }

    return content.Build
}

func (self api) matchProject(project Project) bool {
    return project.Name == "PressReader"
}

func (self api) matchBuildType(buildType BuildType) bool {
    return true;
}

func (self api) processProject(project Project) {
    if(self.matchProject(project)) {

        fmt.Println(project.Name)

        buildTypes := self.getBuildTypes(project)

        done := make(chan bool, len(buildTypes))
        for _, buildType := range buildTypes {
            go func(b BuildType) {
                self.processBuildType(b)
                //fmt.Println(b)
                done<-true
            }(buildType)
        }
        for range buildTypes {
            <-done
        }
        fmt.Println("DONE WITH: " + project.Name)
    }
}

func (self api) processBuildType(buildType BuildType) {
    if(self.matchBuildType(buildType)) {

        builds := self.getBuilds(buildType)

        fmt.Print("len(builds)=")
        fmt.Println(len(builds))

        done := make(chan bool, len(builds))
        for _, build := range builds {
            go func(b Build) {
                fmt.Println(b)
                state := buildState{}
                state.buildTypeId = buildType.Id
                state.currentState = build.Status == SUCCESS

                fmt.Print("build.State=")
                fmt.Print(build.State)
                fmt.Println()

                self.buildStateChanged <- state
                done<-true
            }(build)
        }
        for range builds {
            <-done
        }

        fmt.Println("DONE WITH: " + buildType.Name)
    }
}


func (self api) apiRequest(url string) (*http.Response, error) {
    fmt.Println("---URL: " + url)
    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(self.config.Login, self.config.Password)
    req.Header.Add("Accept", "application/json")

    client := &http.Client{}
    return client.Do(req)
}

func (self api) readJson(url string, obj interface{}) error {
    res, err := self.apiRequest(url)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    bytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }
    fmt.Println("RAW:" + string(bytes))

    if json.Unmarshal(bytes, obj) != nil {
        return err
    }

    return nil
}

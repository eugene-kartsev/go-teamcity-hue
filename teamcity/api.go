package teamcity

import (
    "go-teamcity-hue/config"
    "fmt"
    "errors"
    "go-teamcity-hue/utils"
)

type Api interface {
    GetProjects(loadBuildTypes bool) ([]Project, error)
    GetBuildTypes(projectId *string) ([]BuildType, error)
    GetBuildState(buildTypeId *string) (*Build, error)
}

type api struct {
    config        config.IConfig
}


func Create(cfg config.IConfig) (Api, error) {

    _api := api {
        config: cfg,
    }

    if _, err := _api.GetProjects(false); err != nil {
        fmt.Println(err.Error())
        return nil, errors.New("TeamCity server can't be reached. Check your configuration.")
    }

    return &_api, nil
}

func (self api) GetProjects(loadBuildTypes bool) ([]Project, error) {
    url := self.config.GetTeamcityApiUrl() + "projects"
    fmt.Println("teamcity. Request url: " + url)

    var content projectList
    if err := self.readJson(url, &content); err != nil {
        return nil, err
    }

    if loadBuildTypes && len(content.Projects) > 0 {
        projects := content.Projects
        done     := make(chan bool, len(projects))

        for index, _ := range content.Projects {
            project := &content.Projects[index]
            project.BuildTypes = []BuildType{}

            go (func(p *Project) {
                buildTypes, err := self.GetBuildTypes(&p.Id)

                if err == nil {
                    p.BuildTypes = buildTypes
                }

                done<-true
            })(project);
        }

        for range projects {
            <-done
        }
    }

    return content.Projects, nil
}

func (self api) GetBuildTypes(projectId *string) ([]BuildType, error) {
    url := self.config.GetTeamcityApiUrl() + "projects/id:" + (*projectId)
    fmt.Println("teamcity. Request url: " + url)

    var content projectBuildTypes
    if err := self.readJson(url, &content); err != nil {
        fmt.Println("ERROR: " + err.Error())
        return nil, err
    }

    return content.BuildTypes.BuildTypes, nil
}

func (self api) GetBuildState(buildTypeId *string) (*Build, error) {
    url := self.config.GetTeamcityApiUrl() + "buildTypes/id:"+ (*buildTypeId) +"/builds?count=1&start=0&running=false&canceled=false"
    fmt.Println("teamcity. Request url: " + url)

    var content buildList
    if err := self.readJson(url, &content); err != nil {
        fmt.Println("ERROR: " + err.Error())
        return nil, err
    }

    if len(content.Builds) == 1 {
        return &content.Builds[0], nil
    }

    return nil, nil
}

func (self api) readJson(url string, obj interface{}) error {
    login, password := self.config.GetTeamcityCredentials()

    return utils.ReadJsonWithCredentials(url, obj, login, password)
}

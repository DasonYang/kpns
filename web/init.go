package web

import (
    "os"
    "log"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "kpns/database"

    "gopkg.in/yaml.v2"
)

// var (
//     Permissions = map[string][]string{"all":{"allow", "appkey", "search"}, 
//                                       "editor":{"search",},
//                                       "basic":{"search",}}
// )

type Permission struct {
    Readable    bool    `yaml:"readable"`
    Writable    bool    `yaml:"writable"`

}

var PermissionGroups []string

var Permissions = make(map[string]map[string]Permission)

var (
    TemplatePath string
    dbClient    database.DatabaseClient
)

const (
    db_name = "tpns"  
)

func Init(dbc database.DatabaseClient) error {
    path, _ := os.Getwd()
    TemplatePath = filepath.Join(path, "web/templates")
    fmt.Printf("TemplatePath = %v\n", TemplatePath)
    dbClient = dbc

    yamlPath := filepath.Join(path, "config/permission.yaml")

    yamlFile, err := ioutil.ReadFile(yamlPath)

    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }

    err = yaml.Unmarshal(yamlFile, &Permissions)

    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    for group := range Permissions {
        PermissionGroups = append(PermissionGroups, group)
    }

    // fmt.Printf("Permissions = %v\n", Permissions)

    return nil
}
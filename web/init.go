package web

import (
	"fmt"
	"io/ioutil"
	"kpns/database"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// var (
//     Permissions = map[string][]string{"all":{"allow", "appkey", "search"},
//                                       "editor":{"search",},
//                                       "basic":{"search",}}
// )

type Permission struct {
	Readable bool `yaml:"readable"`
	Writable bool `yaml:"writable"`
}

type templateGenerator interface {
	genInput() map[string]interface{}
}

var permissionGroups []string

var permissions = make(map[string]map[string]Permission)

var (
	templatePath string
	dbClient     database.DatabaseClient
)

const (
	dbName = "tpns"
)

// Init - Initial web function
func Init(dbc database.DatabaseClient) error {
	path, _ := os.Getwd()
	templatePath = filepath.Join(path, "web/templates")
	fmt.Printf("TemplatePath = %v\n", templatePath)
	dbClient = dbc

	yamlPath := filepath.Join(path, "config/permission.yaml")

	yamlFile, err := ioutil.ReadFile(yamlPath)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &permissions)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for group := range permissions {
		permissionGroups = append(permissionGroups, group)
	}

	// fmt.Printf("Permissions = %v\n", Permissions)

	return nil
}

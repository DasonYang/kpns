package config

import (
    "fmt"
    "io/ioutil"
    //"runtime"

    "gopkg.in/yaml.v2"
)

type CfgYaml struct {
    Core SectionCore `yaml:"core"`
    API SectionAPI `yaml:"api"`
}

type SectionCore struct {
    Port string `yaml:"port"`
}

type SectionAPI struct {
    PushURI string `yaml:"push_uri"`
}

func Echo() {
    fmt.Printf("Config.Echo\n")
}

func LoadConfigYaml(cfgPath string) (CfgYaml, error) {
    var config CfgYaml

    cfgFile, err := ioutil.ReadFile(cfgPath)

    if err != nil {
        return config, err
    }

    err = yaml.Unmarshal(cfgFile, &config)

    if err != nil {
        return config, err
    }

    return config, nil
}


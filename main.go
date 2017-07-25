package main


import (
    "fmt"
    "flag"
    "log"

    "kpns/config"
    "kpns/kpns"
    // "kpns/database"
)

var Version = "No Version Provided"

func main(){
    fmt.Printf("Hello kpns\n")

    //opts := config.CfgYaml{}

    var configFile string

    flag.StringVar(&configFile, "c", "", "Configuration File Path")
    flag.Parse()

    config.Echo()

    fmt.Printf("config file = %v\n", configFile)

    var err error

    kpns.Configs = config.BuildDefaultConfig()

    if configFile != "" {
        kpns.Configs, err = config.LoadConfigYaml(configFile)

        if err != nil {
            log.Printf("Load yaml config file error: '%v'", err)

            return
        }
    }

    fmt.Printf("%v\n", kpns.Configs.Core.Port)

    // kpns.DBClient = database.NewDB(kpns.Configs)

    kpns.RunServer()
}
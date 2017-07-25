package config

import (
    "fmt"
    "io/ioutil"
    "runtime"

    "gopkg.in/yaml.v2"
)

type CfgYaml struct {
    Core SectionCore `yaml:"core"`
    API SectionAPI `yaml:"api"`
    DB  SectionDB   `yaml:"db"`
}

type SectionCore struct {
    Domain          string          `yaml:"domain"`
    Port            string          `yaml:"port"`
    WorkerNum       int64           `yaml:"worker_num"`
    QueueNum        int64           `yaml:"queue_num"`
    Database        string          `yaml:"database"`
}

type SectionAPI struct {
    PushURI string `yaml:"push_uri"`
}

type SectionDB struct {
    Domain          string          `yaml:"domain"`
    Port            string          `yaml:"port"`
    Name            string          `yaml:"name"`
}

func Echo() {
    fmt.Printf("Config.Echo\n")
}

// BuildDefaultPushConf is default config setting.
func BuildDefaultConfig() CfgYaml {
    var conf CfgYaml

    // Core
    conf.Core.Port = "localhost"
    conf.Core.Port = "7379"
    conf.Core.WorkerNum = int64(runtime.NumCPU())
    conf.Core.QueueNum = int64(8192)
    // conf.Core.Mode = "release"
    // conf.Core.Sync = false
    // conf.Core.SSL = false
    // conf.Core.CertPath = "cert.pem"
    // conf.Core.KeyPath = "key.pem"
    // conf.Core.MaxNotification = int64(100)
    // conf.Core.HTTPProxy = ""
    // conf.Core.PID.Enabled = false
    // conf.Core.PID.Path = "gorush.pid"
    // conf.Core.PID.Override = false
    // conf.Core.AutoTLS.Enabled = false
    // conf.Core.AutoTLS.Folder = ".cache"
    // conf.Core.AutoTLS.Host = ""

    // Api
    conf.API.PushURI = "/api/push"
    // conf.API.StatGoURI = "/api/stat/go"
    // conf.API.StatAppURI = "/api/stat/app"
    // conf.API.ConfigURI = "/api/config"
    // conf.API.SysStatURI = "/sys/stats"
    // conf.API.MetricURI = "/metrics"

    // Android
    // conf.Android.Enabled = false
    // conf.Android.APIKey = ""
    // conf.Android.MaxRetry = 0

    // iOS
    // conf.Ios.Enabled = false
    // conf.Ios.KeyPath = "key.pem"
    // conf.Ios.Password = ""
    // conf.Ios.Production = false
    // conf.Ios.MaxRetry = 0

    // log
    // conf.Log.Format = "string"
    // conf.Log.AccessLog = "stdout"
    // conf.Log.AccessLevel = "debug"
    // conf.Log.ErrorLog = "stderr"
    // conf.Log.ErrorLevel = "error"
    // conf.Log.HideToken = true

    // conf.Stat.Engine = "memory"
    // conf.Stat.Redis.Addr = "localhost:6379"
    // conf.Stat.Redis.Password = ""
    // conf.Stat.Redis.DB = 0

    // conf.Stat.BoltDB.Path = "bolt.db"
    // conf.Stat.BoltDB.Bucket = "gorush"

    // conf.Stat.BuntDB.Path = "bunt.db"
    // conf.Stat.LevelDB.Path = "level.db"

    return conf
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


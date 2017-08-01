package kpns

import (
    "os"
    "fmt"
    "path/filepath"
    "kpns/config"
    "kpns/database"

    log "github.com/sirupsen/logrus"

)

const ( // Device
    InvalidUID = iota + 1       //1
    UIDLengthError              //2
)

// const ( // Client
//     InvalidAppid = iota + 101
// )

const ( // Mapping
    InvalidAppid = iota + 201
)

const (// DB Result
    InsertDBError = iota + 301
)

type KPNSResult struct {
    Code    string
    Msg     string
}

var (
    Configs config.CfgYaml
    QueueNotification chan PushNotification
    DBClient    database.DatabaseClient
    ErrorCode   string = fmt.Sprintf("E%03d", 0)
    TemplatePath string = getTemplatePath()
)

func getTemplatePath() string {
    fmt.Println("getTemplatePath")
    path, _ := os.Getwd()
    return filepath.Join(path, "kpns/templates")
}


func UIDVerify(uid string) (map[string]interface{}, error) {
    log.WithFields(log.Fields {
        "uid": uid,
    }).Info("UIDVerify")
    data := DBClient.ReadOne("tpns", "allow", map[string]interface{}{"key":uid})

    log.WithFields(log.Fields {
        "data": data,
    }).Info("Returned UID")
    return data, nil
}
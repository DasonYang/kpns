package kpns

import (
    "fmt"
    "errors"
    "kpns/utils"

    log "github.com/sirupsen/logrus"
)

func divreverse(n uint64) uint64 {

    var ret uint64
    for i := n; i > 0; {
        ret = ret<<1 + i%2
        i = i / 2
    }

    return ret
}

type ClientData struct {
    Key     string                  `bson:"key"`
    Value   map[string]interface{}  `bson:"value"`
}


func ClientDo(data map[string]interface{}) error {
    log.WithFields(log.Fields {
        "data": data,
    }).Info("Client Args")
    udid := fmt.Sprintf("%v", data["udid"])
    appid := fmt.Sprintf("%v", data["appid"])
    os := fmt.Sprintf("%v", data["os"])
    key := utils.Sum128toString(udid, appid, os)
    log.WithFields(log.Fields {
        "key": key, 
    }).Info("Client Key")

    if key != "##" {
        var count int
        clientInfo := make(map[string]interface{})
        clientInfo["key"] = key
        

        query := make(map[string]interface{})
        query["key"] = key
        ret := DBClient.ReadOne("tpns", "client", query)

        value := utils.ConvertInterfaceToMap(ret["value"])
        // log.Printf("value = %v\n", value)

        if val, ok := value["count"]; ok{
            count = val.(int)
        }

        count++
        data["count"] = count

        clientInfo["value"] = data

        // log.Printf("Result of client with key : %v is %v\n", key, ret)
        err := DBClient.Write("tpns", "client", clientInfo)
        if err != nil {
            log.Debug("Insert db error.")
            return errors.New("Insert db error.")
        }
        return nil
    }
    return errors.New("Making key error")
}

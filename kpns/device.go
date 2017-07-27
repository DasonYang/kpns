package kpns

import (
    "fmt"
    "time"
    // "errors"

    "kpns/utils"

    log "github.com/sirupsen/logrus"
)

func DeviceDo(data map[string]interface{}) KPNSResult {
    log.WithFields(log.Fields {
        "data" : data,
    }).Info("Device Args")

    var uid string

    if val, ok := data["uid"]; ok && len(val.(string)) == 20 {
        log.WithFields(log.Fields {
            "uid" : val,
        }).Info("UID")
        uid = fmt.Sprintf("%v", val)

        if len(uid) != 20 {
            // return errors.New("UID length error")
            return KPNSResult{ErrorCode + fmt.Sprintf("%03d", UIDLengthError), "UID length error"}
        }

        result, _ := UIDVerify(uid)

        log.WithFields(log.Fields {
            "ret" : result,
        }).Info("UIDVerify")

        var count int
        info := make(map[string]interface{})
        info["key"] = uid
        
        query := make(map[string]interface{})
        query["key"] = uid
        ret := DBClient.Read("tpns", "device", query)

        value := utils.ConvertInterfaceToMap(ret["value"])

        if val, ok := value["count"]; ok{
            count = val.(int)
        }

        count++
        data["count"] = count

        if val, ok := value["first_time"]; ok {
            log.Info("Get first_time")
            data["first_time"] = val
            data["update"] = int32(time.Now().Unix())
        } else {
            data["first_time"] = int32(time.Now().Unix())
        }

        log.WithFields(log.Fields {
            "data" : data,
        }).Info("DATA")

        info["value"] = data

        // log.Printf("Result of client with key : %v is %v\n", key, ret)
        err := DBClient.Write("tpns", "device", info)
        if err != nil {
            log.Debug("Insert db error.")
            return KPNSResult{ErrorCode + fmt.Sprintf("%03d", InsertDBError), "Insert DB Error."}
        }
    } else {
        // return errors.New("UID length error")
        return KPNSResult{ErrorCode + fmt.Sprintf("%03d", UIDLengthError), "UID length error"}
    }

    return KPNSResult{"200", "Success"}
}
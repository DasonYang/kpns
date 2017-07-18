package kpns

import (
    "log"
    "bytes"
    "github.com/reusee/mmh3"
)

type ClientData struct {
    Key     string                  `bson:"key"`
    Value   map[string]interface{}  `bson:"value"`
}

func GetClientKey(udid string, appid string, platform string) string {
    var key string
    var buffer bytes.Buffer
    buffer.WriteString(udid)
    buffer.WriteString(appid)
    buffer.WriteString(platform)
    key = string(mmh3.Sum128(buffer.Bytes()))

    log.Printf("hashed key = %v\n", key)

    return key
}
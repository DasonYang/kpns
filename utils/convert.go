package utils

import (
    "log"
    "reflect"
)


func ConvertInterfaceToMap(in interface{}) map[string]interface{} {
    tp := reflect.ValueOf(in)
    log.Printf("kind = %v\n", tp.Kind())
    if tp.Kind() == reflect.Map {
        data := make(map[string]interface{})
        for _, key := range tp.MapKeys() {
            value := tp.MapIndex(key)
            data[key.String()] = value.Interface()
        }
        return data
    }
    return nil
}
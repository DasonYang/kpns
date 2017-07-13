package kpns

import (
    "testing"
    "encoding/json"
    "fmt"
    "reflect"
)

type TestStruct struct {
    Name    string      `json:"name"`
    Id      int         `json:"id"`
}

func TestSendToQueue(t *testing.T) {
    // var count int

    req := PushRequest{}

    // str := `{"notifications" :[{"tokens" : ["token_1", "token_2"],"platform" : 1},{"tokens" : ["token_3", "token_4"],"platform" : 2}]}`
    str := []byte(`{"notifications" :[{"tokens" : ["token_1", "token_2"],"platform" : 1},{"tokens" : ["token_3", "token_4"],"platform" : 2}]}`)
    fmt.Println(reflect.TypeOf(str), str)
    // str := `{"notifications" :[{"platform" : 1},{"tokens" : ["token_3", "token_4"],"platform" : 2}]}`


    err := json.Unmarshal(str, &req)

    if err != nil {
        fmt.Printf("Error = %v\n", err)
    }

    for _, notification := range req.Notifications {
        fmt.Printf("notification = %v, tokens = %v, platform = %v\n", notification, notification.Tokens, notification.Platform)
    }



    // count = queueNotification(req)

    // fmt.Printf("Total count of notifications = %v\n", count)
}

func TestUnmarshal(t *testing.T) {

    req := TestStruct{}
    str := []byte(`{"name" : "JerryYang", "id":1919}`)

    err := json.Unmarshal(str, &req)

    if err != nil {
        fmt.Printf("Error = %v\n", err)
    }
    fmt.Printf("name = %v, id = %v\n", req.Name, req.Id)

}
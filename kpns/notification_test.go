package kpns

import (
    "testing"
    "encoding/json"
    "fmt"
    "reflect"
    "log"

    "github.com/sideshow/apns2"
    "github.com/sideshow/apns2/certificate"
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

func TestPushIos(t *testing.T) {
    cert, err := certificate.FromPemFile("./ios.pem", "")
    if err != nil {
        log.Fatal("Cert Error:", err)
    }

    notification := &apns2.Notification{}
    notification.DeviceToken = "991094614d098364a30f4ea448e17ae39f16013866ebcf18af4b6a540fa36009"
    notification.Topic = "com.tutk.p2pcamlive.2"
    notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

    client := apns2.NewClient(cert).Production()
    res, err := client.Push(notification)

    if err != nil {
        log.Fatal("Error:", err)
    }

    fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)

    if res.Sent() {
        log.Println("Sent:", res.ApnsID)
    } else {
        fmt.Printf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
    }

}
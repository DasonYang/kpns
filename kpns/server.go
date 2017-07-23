package kpns

import (
    "fmt"
    "log"
    // "reflect"
    "kpns/utils"
    "encoding/json"


    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)

func pushHandler(c *gin.Context) {
    fmt.Printf("Push Handler")
    var data PushRequest
    var msg string

    if err := c.BindJSON(&data); err != nil {
        msg = "Missing notifications field."
        // LogAccess.Debug(msg)
        // abortWithError(c, http.StatusBadRequest, msg)
        return
    }

    if len(data.Notifications) == 0 {
        msg = "Notifications field is empty."
        // LogAccess.Debug(msg)
        // abortWithError(c, http.StatusBadRequest, msg)
        return
    }

    for _, notification := range data.Notifications {
        fmt.Printf("notification = %v, tokens = %v, platform = %v, msg = %v\n", notification, notification.Tokens, notification.Platform, msg)
    }

    // if int64(len(form.Notifications)) > PushConf.Core.MaxNotification {
    //     msg = fmt.Sprintf("Number of notifications(%d) over limit(%d)", len(form.Notifications), PushConf.Core.MaxNotification)
    //     LogAccess.Debug(msg)
    //     abortWithError(c, http.StatusBadRequest, msg)
    //     return
    // }

    // counts, logs := queueNotification(form)

    // c.JSON(http.StatusOK, gin.H{
    //     "success": "ok",
    //     "counts":  counts,
    //     "logs":    logs,
    // })
}

func pushHandlerGET(c *gin.Context) {
    query := c.Request.URL.Query()

    cmd := query["cmd"][0]
    log.Printf("query = %v\n", query)
    delete(query, "cmd")

    client := make(map[string]interface{})
    client["key"] = "1029384756"
    values := make(map[string]interface{})

    for k, v := range query{
       values[k] = v[0]
    }
    client["value"] = values
    fmt.Printf("client = %v\n", client)

    jsonString, err := json.Marshal(query)

    if err != nil {
        panic(err)
    }

    log.Printf("json string = ", jsonString)

    var dat map[string]interface{}

    err = json.Unmarshal(jsonString, &dat)

    log.Printf("dat = ", dat)

    session, err := mgo.Dial("localhost")

    if err != nil {
        panic(err)
    }

    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    

    switch cmd{
    case "hello":
        c.String(200, "Hi")
    case "event":
        fallthrough
    case "raise_event":
        c.String(200, "Event")
    case "mapping":
        fallthrough
    case "reg_mapping":
        c.String(200, "Mapping")
    case "rm_mapping":
        fallthrough
    case "unreg_mapping":
        c.String(200, "Unmapping")
    case "client":
        fallthrough
    case "reg_client":
        // collection := session.DB("test").C("client")
        // err = collection.Insert(client)
        // if err != nil {
        //     panic(err)
        // }
        log.Printf("udid : %v, appid : %v, os : %v\n", query["udid"][0], query["appid"][0], query["os"][0])
        key := utils.Sum128toString(query["udid"][0], query["appid"][0], query["os"][0])
        log.Printf("Get key : %v\n", key)
        c.String(200, "Client")
    default:
        c.String(200, "Error")
    }
    
}

func RunServer() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    r.POST("/api/push", pushHandler)
    r.GET("/tpns", pushHandlerGET)
    r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
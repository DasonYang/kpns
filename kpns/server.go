package kpns

import (
    "fmt"
    "github.com/gin-gonic/gin"
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
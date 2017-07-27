package kpns

import (
    "github.com/gin-gonic/gin"

    log "github.com/sirupsen/logrus"
)

func kpnsGETHandler(c *gin.Context) {
    query := c.Request.URL.Query()

    cmd := query["cmd"][0]
    log.WithFields(log.Fields{"query" : query}).Info()
    delete(query, "cmd")

    data := make(map[string]interface{})

    for k, v := range query{
       data[k] = v[0]
    }

    log.Printf("data = %v\n", data)

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
        err := ClientDo(data)

        if err != nil {
            panic(err)
        } else {
            c.String(200, "Success")
        }
    case "device":
        fallthrough
    case "reg_server":
        result := DeviceDo(data)

        c.String(200, result.Code+ " " +result.Msg)
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

    // r.POST("/api/push", pushHandler)
    r.GET("/tpns", kpnsGETHandler)
    r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
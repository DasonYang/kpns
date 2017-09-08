package kpns

import (
    // "github.com/gin-gonic/gin"

    // log "github.com/sirupsen/logrus"
    "fmt"
    "net/http"
    "encoding/json" 

    "kpns/web"
)

// func kpnsGETHandler(c *gin.Context) {
//     query := c.Request.URL.Query()

//     cmd := query["cmd"][0]
//     log.WithFields(log.Fields{"query" : query}).Info()
//     delete(query, "cmd")

//     data := make(map[string]interface{})

//     for k, v := range query{
//        data[k] = v[0]
//     }

//     log.Printf("data = %v\n", data)

//     switch cmd{
//     case "hello":
//         c.String(200, "Hi")
//     case "event":
//         fallthrough
//     case "raise_event":
//         c.String(200, "Event")
//     case "mapping":
//         fallthrough
//     case "reg_mapping":
//         c.String(200, "Mapping")
//     case "rm_mapping":
//         fallthrough
//     case "unreg_mapping":
//         c.String(200, "Unmapping")
//     case "client":
//         fallthrough
//     case "reg_client":
//         err := ClientDo(data)

//         if err != nil {
//             panic(err)
//         } else {
//             c.String(200, "Success")
//         }
//     case "device":
//         fallthrough
//     case "reg_server":
//         result := DeviceDo(data)

//         c.String(200, result.Code+ " " +result.Msg)
//     default:
//         c.String(200, "Error")
//     }
// }

// func RunServer() {
//     r := gin.Default()
//     r.GET("/ping", func(c *gin.Context) {
//         c.JSON(200, gin.H{
//             "message": "pong",
//         })
//     })

//     // r.POST("/api/push", pushHandler)
//     r.GET("/tpns", kpnsGETHandler)
//     r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
// }

func kpnsHandler(w http.ResponseWriter, r *http.Request) {

    var data = make(map[string]interface{})

    err := r.ParseForm()
    if err != nil {
        panic(err)
    }

    if r.Method == http.MethodPost {
        content_type := r.Header["Content-Type"][0]

        switch content_type {
        case "application/json":
            if r.Body != nil {
                err = json.NewDecoder(r.Body).Decode(&data)

                if err != nil {
                    panic(err)
                }
                defer r.Body.Close()
            }
        case "application/x-www-form-urlencoded":
            query := r.PostForm
            for k, v := range query {
                data[k] = v[0]
            }
        default:
            r.ParseMultipartForm(0)
            fmt.Println("Not support ", r.PostForm)
        }
    } else if r.Method == http.MethodGet {
        query := r.Form
        for k, v := range query {
            data[k] = v[0]
        }
    }

    fmt.Printf("Data =  %v\n", data)
    cmd := data["cmd"]
    delete(data, "cmd")

    switch cmd {
    case "hello":
        fmt.Fprintln(w, "200 Success")
    case "event":
        fallthrough
    case "raise_event":
        fmt.Fprintln(w, "200 Success")
    case "mapping":
        fallthrough
    case "reg_mapping":
        fmt.Fprintln(w, "200 Success")
    case "rm_mapping":
        fallthrough
    case "unreg_mapping":
        fmt.Fprintln(w, "200 Success")
    case "client":
        fallthrough
    case "reg_client":
        err := ClientDo(data)

        if err != nil {
            panic(err)
        }
        fmt.Fprintln(w, "200 Success")
    case "device":
        fallthrough
    case "reg_server":
        result := DeviceDo(data)

        // c.String(200, result.Code+ " " +result.Msg)
        fmt.Fprintln(w, "200 Success ", result)
    default:
        fmt.Fprintln(w, "200 Success")
    }

}

func RunServer() {
    err := web.Init(DBClient)

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/tpns", kpnsHandler)
    http.HandleFunc("/login", web.LoginHandler)
    http.HandleFunc("/logout", web.LogoutHandler)
    http.Handle("/allow", web.AuthMiddleware(http.HandlerFunc(web.AllowHandler)))
    http.Handle("/search", web.AuthMiddleware(http.HandlerFunc(web.SearchHandler)))
    http.Handle("/appkey", web.AuthMiddleware(http.HandlerFunc(web.AppKeyHandler)))
    http.Handle("/lang", web.AuthMiddleware(http.HandlerFunc(web.LangHandler)))
    http.Handle("/account", web.AuthMiddleware(http.HandlerFunc(web.AccountHandler)))
    http.Handle("/log", web.AuthMiddleware(http.HandlerFunc(web.LogHandler)))

    http.ListenAndServe(":8080", nil)
}
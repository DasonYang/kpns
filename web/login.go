package web

import (
    "io"
    "fmt"
    "time"
    "strconv"
    "net/http"
    "math/rand"
    "crypto/md5"
    "encoding/hex"
    "html/template"
    "encoding/base64"


    "kpns/utils"
)

const (
    MessageType_Error = iota
    MessageType_Warning
)

func checkAccount(username, password string) (string, string) {
    fmt.Printf("username = %v, password = %v\n", username, password)

    h := md5.New()
    io.WriteString(h, username)
    io.WriteString(h, password)
    realPassword := hex.EncodeToString(h.Sum(nil))

    h = md5.New()
    io.WriteString(h, username)
    dpwd, _ := base64.StdEncoding.DecodeString(password)
    io.WriteString(h, string(dpwd))
    testPassword := hex.EncodeToString(h.Sum(nil))

    fmt.Printf("realPassword = %v, testPassword = %v\n", realPassword, testPassword)

    ret := dbClient.ReadOne("tpns", "account", map[string]interface{}{"key":username})

    fmt.Printf("ret = %v, ret = %v\n", ret, ret)

    if val, ok := ret["value"]; ok {

        value := val.(map[string]interface{})

        // Make sure if pswd is string
        if pwd, ok := value["pswd"].(string); ok {
            if pwd == realPassword || pwd == testPassword { //Password confirm 
                value["last"] = int32(time.Now().Unix())
                if val, ok := value["count"]; ok {
                    value["count"] = val.(int) + 1
                } else {
                    value["count"] = 1
                }
                rand.Seed(time.Now().UnixNano())
                r := strconv.FormatInt(rand.Int63(), 10)
                token := utils.Mmh3py(r)
                value["token"] = token
                fmt.Printf("token = %v\n", token)
                var mode string
                // if val, ok := value["mode"].(string); ok {
                //     mode = val
                // }
                ret["value"] = value
                dbClient.Write("tpns", "account", ret)

                return token, mode
            }
        }
    }



    fmt.Printf("ret = %v\n", ret)

    return "", ""
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("================================Login=================================")
    if r.Method == "GET" {
        fmt.Println("================================Login.GET=================================")
        // info := make(map[string]interface{})
        // msg :=  r.Cookie("msg")

        // info := make(map[string]interface{})
            
        // if msg != "" {
        //     info["Type"] = MessageType_Warning
        //     info["Message"] = msg
        // }
        t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
        if err != nil {
            fmt.Printf("Error = %v\n", err)
            panic(err)
        }
        t.Execute(w, nil)
    } else {
        fmt.Println("================================Login.POST=================================")
        r.ParseMultipartForm(0)
        // logic part of log in
        username := r.PostForm["user"][0]
        password := r.PostForm["pswd"][0]
        // url = r.Cookie("url")

        token, mode := checkAccount(username, password)

        if token == ""  {
            // user := r.Cookie("user")
            
            info := make(map[string]interface{})
            
            // if user != "" && msg == "" {
            //     info["Type"] = MessageType_Error
            //     info["Message"] = "Login Error!"
            // } else if msg != "" {
            //     info["Type"] = MessageType_Warning
            //     info["Message"] = msg
            // }

            t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
            if err != nil {
                fmt.Printf("Error = %v\n", err)
                panic(err)
            }
            t.Execute(w, info)
        } else {
            // Login success

            expiration := time.Now()
            expiration = expiration.Add(time.Minute * time.Duration(1))
            fmt.Printf("expiration = %v, token = %v, mode = %v\n", expiration, token, mode)
            // cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
            http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: expiration})
            http.SetCookie(w, &http.Cookie{Name: "mode", Value: mode, Expires: expiration})
            http.SetCookie(w, &http.Cookie{Name: "user", Value: username, Expires: expiration})
            http.Redirect(w, r, "/allow", http.StatusSeeOther)
        }
    }
}
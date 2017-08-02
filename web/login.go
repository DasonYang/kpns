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
                last := time.Now().Unix()
                value["last"] = last
                var count int
                if val, ok := value["count"]; ok {
                    count = val.(int)
                }
                count++
                value["count"] = count
                rand.Seed(time.Now().UnixNano())
                r := strconv.FormatInt(rand.Int63(), 10)
                token := utils.Mmh3py(r)
                value["token"] = token
                fmt.Printf("token = %v\n", token)
                var mode string
                if val, ok := value["mode"].(string); ok {
                    mode = val
                }

                return token, mode
            }
        }
    }



    fmt.Printf("ret = %v\n", ret)

    return "", ""
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintln(w, "200 Success"+" "+r.URL.Path[1:])
    if r.Method == "GET" {
        // info := make(map[string]interface{})
        t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
        if err != nil {
            fmt.Printf("Error = %v\n", err)
            panic(err)
        }
        t.Execute(w, nil)
    } else {
        r.ParseMultipartForm(0)
        // logic part of log in
        username := r.PostForm["user"][0]
        password := r.PostForm["pswd"][0]

        token, mode := checkAccount(username, password)

        if token == ""  {
            fmt.Println("username is empty")
            info := make(map[string]interface{})
            info["Error"] = true
            info["Message"] = "Login Error!"

            t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
            if err != nil {
                fmt.Printf("Error = %v\n", err)
                panic(err)
            }
            t.Execute(w, info)
        } else {
            // r.SetBasicAuth("jerryyang", "123")
            // username, password, ok := r.BasicAuth()
            expiration := time.Now()
            expiration = expiration.AddDate(1, 0, 0)
            fmt.Printf("token = %v, mode = %v\n", token, mode)
            // cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
            http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: expiration})
            http.SetCookie(w, &http.Cookie{Name: "mode", Value: mode, Expires: expiration})
            http.SetCookie(w, &http.Cookie{Name: "user", Value: username, Expires: expiration})
            http.Redirect(w, r, "/allow", http.StatusSeeOther)
        }
    }
}
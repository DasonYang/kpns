package web

import(
    "fmt"
    // "strconv"
    "net/http"
    "html/template"
)

type LogData struct {
    DateTime        string
    Status          string
    Content         string

}

func LogHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println("================================Account=================================", r.Context().Value("Writable"))
    // var writable = r.Context().Value("Writable").(bool)
    // var query = make(map[string]interface{})
    // var pageIdx, limit int
    // // var account string
    // var resultMsg string
    // var popAnalog bool
    // getInput
    // genInput := func(limit, page int, query map[string]interface{}, success bool, msg string, writable bool) map[string]interface{} {
    //     //fmt.Printf("limit = %v, page = %v, note = %v, query = %v\n", limit, page, note, query)
    //     var input = make(map[string]interface{})
    //     var params = make(map[string]interface{})
    //     var accountList []AccountData
    //     pageIdx := page
    //     displayLimit := limit

    //     if limit < 20 {displayLimit = 20}
    //     if page == 0 {pageIdx = 1}

    //     params["skip"] = (pageIdx-1)*displayLimit
    //     params["limit"] = displayLimit

    //     // fmt.Printf("params = %v\n", params)

    //     qs, count := dbClient.ReadAll(db_name, "account", query, params)

    //     for _, allow := range qs {
    //         var data AccountData
    //         if str, f := allow["key"].(string); f{data.User = str}

    //         value := allow["value"].(map[string]interface{})

    //         if str, f := value["last"].(string); f {
    //             data.LastTime = str
    //         } else if ts, f := value["last"].(float64); f {
    //             tm := time.Unix(int64(ts), 0)
    //             data.LastTime = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
    //         } else if ts, f := value["last"].(int); f {
    //             tm := time.Unix(int64(ts), 0)
    //             data.LastTime = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
    //         }
            
    //         if str, f := value["first"].(string); f {
    //             data.FirstTime = str
    //         } else if ts, f := value["first"].(float64); f {
    //             tm := time.Unix(int64(ts), 0)
    //             data.FirstTime = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
    //         } else if ts, f := value["first"].(int); f {
    //             tm := time.Unix(int64(ts), 0)
    //             data.FirstTime = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
    //         }
            
    //         if str, f := value["mode"].(string); f {data.Mode = str}
    //         if cnt, f := value["count"].(int); f {data.Count = cnt}
            
    //         // fmt.Printf("type of limit = %v\n", reflect.TypeOf(value["limit"]))
    //         accountList = append(accountList, data)
    //     }

    //     // fmt.Printf("count = %v\n", count)

    //     fmt.Printf("PermissionGroups = %v\n", PermissionGroups)

    //     input["Data"] = accountList
    //     input["Page"] = pageIdx
    //     input["Count"] = count
    //     input["Limit"] = displayLimit
    //     input["Success"] = success
    //     input["Writable"] = writable
    //     input["User"] = ""
    //     input["Pswd"] = ""
    //     input["Perms"] = PermissionGroups
    //     input["Msg"] = msg

    //     if pageIdx > 1 {
    //         input["HasPre"] = true
    //         input["Pre"] = pageIdx-1
    //     }
    //     if (pageIdx * displayLimit) < count {
    //         input["HasNext"] = true
    //         input["Next"] = pageIdx + 1
    //     }

    //     return input
    // }
    // genInput
    t, err := template.ParseFiles(TemplatePath+"/log.tmpl")
    if err != nil {
        fmt.Printf("Error = %v\n", err)
        panic(err)
    }

    if r.Method == "GET" {
        // fmt.Println("================================Account.GET=================================")
        // var active string
        // for key, value := range r.URL.Query() {
        //     fmt.Printf("key = %v, value = %v\n", key, value)
        //     switch key {
        //     case "page":
        //         pageIdx, _ = strconv.Atoi(value[0])
        //     case "limit":
        //         limit, _ = strconv.Atoi(value[0])
        //     case "active":
        //         active = value[0]
        //     case "account":
        //         account = value[0]
        //     }
        // }

        // if active == "del" {
        //     err := dbClient.Delete(db_name, "account", map[string]interface{}{"key":account})
        //     if err != nil {
        //         log.Println(err)
        //     }
        // }

        // input := genInput(limit, pageIdx, query, popAnalog, resultMsg, writable)

        t.Execute(w, nil)
    } else {
        fmt.Println("================================Log.POST=================================")
        r.ParseMultipartForm(0)

        // var mode string
        // var pswd string
        // for key, value := range r.PostForm {
        //     fmt.Printf("key = %v, value = %v\n", key, value)
        //     switch key {
        //     case "limit":
        //         limit, _ = strconv.Atoi(value[0])
        //     case "account":
        //         account = strings.TrimSpace(value[0])
        //     case "secure":
        //         pswd = strings.TrimSpace(value[0])
        //         fmt.Printf("pswd = %v\n", pswd)
        //     case "mode":
        //         mode = strings.TrimSpace(value[0])

        //         if _, ok := Permissions[mode]; !ok {
        //             mode = "default"
        //         }
        //         fmt.Printf("mode = %v\n", mode)
        //     }
        // }
        
        // if _, ok := r.PostForm["search"]; ok {// Search clicked
        //     if limit < 20 {limit = 20}
        //     if pageIdx == 0 {pageIdx = 1}

        //     if len(account) > 0 {
        //         query["key"] = map[string]interface{}{"$regex":account}
        //     }
        // } else if _, ok := r.PostForm["save"]; ok {// Save clicked

        //     if len(account) > 0 && len(pswd) > 0 {

        //         count := dbClient.Count(db_name, "account", map[string]interface{}{"key":account})

        //         if count == 0 {
        //             var encPswd, savedPasswd string
        //             query["key"] = account

        //             encPswd = base64.StdEncoding.EncodeToString([]byte(pswd))

        //             genPswd := func(account, password string) string {
        //                 h := md5.New()
        //                 io.WriteString(h, account)
        //                 io.WriteString(h, password)
        //                 return hex.EncodeToString(h.Sum(nil))
        //             }

        //             if strings.Contains(account, "@") {
        //                 realPassword := genPswd(account, strconv.Itoa(int(time.Now().Unix())))
        //                 fmt.Printf("realPassword = %v\n", realPassword)
        //                 encPswd = base64.StdEncoding.EncodeToString([]byte(realPassword))
        //                 {

        //                     auth := utils.NewLoginAuth(KALAY_USER, KALAY_PASSWORD)

        //                     to := []string{account}
        //                     msg := []byte(fmt.Sprintf(mail_body, account, realPassword))
        //                     err := utils.SendMail("outlook-apacsouth.office365.com:587", auth, KALAY_USER, to, msg)
        //                     if err != nil {
        //                         log.Println("with err:", err)
        //                     }
        //                     fmt.Println("please check mailbox")
        //                 }

        //             } else {
        //                 encPswd = base64.StdEncoding.EncodeToString([]byte(pswd))
        //             }

        //             savedPasswd = genPswd(account, encPswd)
        //             fmt.Printf("savedPasswd = %v\n", savedPasswd)
        //             info := make(map[string]interface{})
        //             info["key"] = account
        //             data := make(map[string]interface{})
        //             data["pswd"] = savedPasswd
        //             data["mode"] = mode
        //             data["first"] = int32(time.Now().Unix())
        //             info["value"] = data

        //             dbClient.Write(db_name, "account", info)

                    

        //             query["key"] = account
        //         } else {
        //             popAnalog = true
        //             resultMsg = "Account exists."
        //         }

        //         // info := map[string]interface{}{"key":account}
        //         // data := make(map[string]interface{})
        //         // data["first"] = time.Now().Format("2006-01-02 15:04:05")
        //         // data["mode"] = mode
        //         // info["pswd"] = data
        //         // dbClient.Write(db_name, "account", info)
        //     }
        // }

        // input := genInput(limit, pageIdx, query, popAnalog, resultMsg, writable)

        t.Execute(w, nil)
    }
}
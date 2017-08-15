package web

import (
    "fmt"
    "time"
    "reflect"
    "strconv"
    "net/http"
    // "io/ioutil"
    "html/template"
)

type AllowData struct {
    UID     string
    Updated string
    Limit   string
    Note    string
}

func genInput(limit, page int, note string, query map[string]interface{}) map[string]interface{} {
    fmt.Printf("limit = %v, page = %v, note = %v, query = %v\n", limit, page, note, query)
    var input = make(map[string]interface{})
    var params = make(map[string]interface{})
    var allowList []AllowData
    pageIdx := page
    displayLimit := limit

    if limit < 20 {displayLimit = 20}
    if page == 0 {pageIdx = 1}

    params["skip"] = (pageIdx-1)*displayLimit
    params["limit"] = displayLimit

    // fmt.Printf("params = %v\n", params)

    qs, count := dbClient.ReadAll("tpns", "allow", query, params)

    for _, allow := range qs {
        var data AllowData
        if str, f := allow["key"].(string); f{data.UID = str}

        value := allow["value"].(map[string]interface{})

        if str, f := value["update_time"].(string); f {data.Updated = str}
        
        if ts, f := value["limit"].(float64); f {
            tm := time.Unix(int64(ts), 0)
            data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
        }

        if ts, f := value["limit"].(int); f {
            tm := time.Unix(int64(ts), 0)
            data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
        }
        
        if str, f := value["note"].(string); f {data.Note = str}
        
        // fmt.Printf("type of limit = %v\n", reflect.TypeOf(value["limit"]))
        allowList = append(allowList, data)
    }

    // fmt.Printf("count = %v\n", count)

    input["Data"] = allowList
    input["Page"] = pageIdx
    input["Count"] = count
    input["Limit"] = displayLimit
    input["HasNote"] = true
    input["Note"] = note

    if pageIdx > 1 {
        input["HasPre"] = true
        input["Pre"] = pageIdx-1
    }
    if (pageIdx * displayLimit) < count {
        input["HasNext"] = true
        input["Next"] = pageIdx + 1
    }

    return input
}

func genData(qs []map[string]interface{}) []AllowData {

    var allowList []AllowData

    for _, allow := range qs {
        var data AllowData
        if str, f := allow["key"].(string); f{data.UID = str}

        value := allow["value"].(map[string]interface{})

        if str, f := value["update_time"].(string); f {data.Updated = str}
        
        if ts, f := value["limit"].(float64); f {
            tm := time.Unix(int64(ts), 0)
            data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
        }
        
        if str, f := value["note"].(string); f {data.Note = str}
        
        // fmt.Printf("data = %v\n", data)
        allowList = append(allowList, data)
    }

    return allowList
}

func AllowHandler(w http.ResponseWriter, r *http.Request) {

    // username, password, ok := r.BasicAuth()
    // fmt.Println("================================Allow=================================")

    var query = make(map[string]interface{})
    var pageIdx, limit int
    var active, note, uid string

    t, err := template.ParseFiles(TemplatePath+"/allow.tmpl")
    if err != nil {
        fmt.Printf("Error = %v\n", err)
        panic(err)
    }

    if r.Method == "GET" {
        fmt.Println("================================Allow.GET=================================")
        
        
        for key, value := range r.URL.Query() {
            fmt.Printf("key = %v, value = %v\n", key, value)
            switch key {
            case "page":
                pageIdx, _ = strconv.Atoi(value[0])
            case "limit":
                limit, _ = strconv.Atoi(value[0])
            case "active":
                active = value[0]
            case "uid":
                uid = value[0]
            case "note":
                note = value[0]
                query["value.note"] = map[string]interface{}{"$regex":note}
            }
        }

        if active == "del" && len(uid) == 20 {
            err := dbClient.Delete("tpns", "allow", map[string]interface{}{"key":uid})
            if err != nil {
                panic(nil)
            }
        }

        input := genInput(limit, pageIdx, note, query)

        t.Execute(w, input)
    } else {
        fmt.Println("================================Allow.POST=================================")
        r.ParseMultipartForm(0)

        var ltime string
        var mode string

        file, _, err := r.FormFile("File")
        if err != nil {
            fmt.Println(err)

        } 
        if file != nil {
            defer file.Close()
        }

        for key, value := range r.PostForm {
            fmt.Printf("key = %v, value = %v\n", key, value)
            switch key {
            case "ltime":
                ltime = value[0]
            case "active":
                active = value[0]
            case "limit":
                limit, _ = strconv.Atoi(value[0])
            case "note":
                note = value[0]
                query["value.note"] = map[string]interface{}{"$regex":note}
            case "File":
            case "mode":
                mode = value[0]
            case "uid":
                uid = value[0]

            }
        }
        

        if active == "search" {
            if limit < 20 {limit = 20}
            if pageIdx == 0 {pageIdx = 1}

            if len(uid) > 0 {
                query["key"] = uid
            }
        } else if active == "save" {

            if len(uid) > 0 && len(ltime) > 0 && len(note) > 0 {
                tm, err := time.Parse("2006/01/02", ltime)

                if err != nil {
                    panic(err)
                }
                fmt.Println(tm.UnixNano(), ltime, mode)

                query["key"] = uid

                info := map[string]interface{}{"key":uid}
                data := make(map[string]interface{})
                data["limit"] = int32(tm.Unix())
                data["note"] = note
                data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
                info["value"] = data
                dbClient.Write("tpns", "allow", info)
                note = "" 
            } else {
                note = ""
            }
        }

        input := genInput(limit, pageIdx, note, query)

        t.Execute(w, input)
    }
}
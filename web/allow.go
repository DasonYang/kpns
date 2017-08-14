package web

import (
    "fmt"
    "time"
    // "reflect"
    "strconv"
    "net/http"
    "html/template"
)

type AllowData struct {
    UID     string
    Updated string
    Limit   string
    Note    string
}

func AllowHandler(w http.ResponseWriter, r *http.Request) {

    // username, password, ok := r.BasicAuth()
    fmt.Println("================================Allow=================================")

    if r.Method == "GET" {
        fmt.Println("================================Allow.GET=================================")
        var params = make(map[string]interface{})
        var allowList []AllowData
        var pageIdx, limit int

        count := dbClient.Count("tpns", "allow", nil)
        
        for key, value := range r.URL.Query() {
            fmt.Printf("key = %v, value = %v\n", key, value)
            switch key {
            case "page":
                pageIdx, _ = strconv.Atoi(value[0])
            case "limit":
                limit, _ = strconv.Atoi(value[0])
            }
        }

        if limit < 20 {
            limit = 20
        }

        if pageIdx == 0 {
            pageIdx = 1
        }
        

        fmt.Printf("pageidx = %v, limit = %v\n", pageIdx, limit)
        params["skip"] = (pageIdx-1)*limit
        params["limit"] = limit

        query := dbClient.ReadAll("tpns", "allow", nil, params)

        fmt.Printf("count = %v\n", len(query))

        for _, allow := range query {
            var data AllowData
            if str, f := allow["key"].(string); f{
                data.UID = str
            }
            value := allow["value"].(map[string]interface{})

            if str, f := value["update_time"].(string); f {
                data.Updated = str
            }
            
            if ts, f := value["limit"].(float64); f {
                tm := time.Unix(int64(ts), 0)
                data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
            }
            
            if str, f := value["note"].(string); f {
                data.Note = str
            }
            
            // fmt.Printf("data = %v\n", data)
            allowList = append(allowList, data)
        }

        fmt.Println("AllowHandler")
        t, err := template.ParseFiles(TemplatePath+"/allow.tmpl")
        if err != nil {
            fmt.Printf("Error = %v\n", err)
            panic(err)
        }
        var input = make(map[string]interface{})
        input["Data"] = allowList
        input["Page"] = pageIdx
        input["Count"] = count
        input["Limit"] = limit

        if pageIdx > 1 {
            input["HasPre"] = true
            input["Pre"] = pageIdx-1
        }
        if (pageIdx * limit) < count {
            input["HasNext"] = true
            input["Next"] = pageIdx + 1
        }

        t.Execute(w, input)
    } else {
        fmt.Println("================================Allow.POST=================================")
        r.ParseMultipartForm(0)
        // logic part of log in
    }
}
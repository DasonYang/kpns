package kpns

import (
    "fmt"
    "time"
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

    

    if r.Method == "GET" {

        var params = make(map[string]interface{})
        var allowList []AllowData
        params["skip"] = 0
        params["limit"] = 20

        query := DBClient.ReadAll("tpns", "allow", params)

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
            
            fmt.Printf("data = %v\n", data)
            allowList = append(allowList, data)
        }

        fmt.Println("AllowHandler")
        t, err := template.ParseFiles(TemplatePath+"/allow.tmpl")
        if err != nil {
            fmt.Printf("Error = %v\n", err)
            panic(err)
        }
        t.Execute(w, allowList)
    } else {
        r.ParseMultipartForm(0)
        // logic part of log in
    }
}
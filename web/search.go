package web

import(
    "fmt"
    "strconv"
    "net/http"
    "html/template"
)

type SearchData struct {
    Key         string
    DB          string
    Value       string

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    var input = make(map[string]interface{})

    t, err := template.ParseFiles(TemplatePath+"/search.tmpl")
    if err != nil {
        fmt.Printf("Error = %v\n", err)
        panic(err)
    }

    if r.Method == "GET" {
        var event, clientno, uid string
        var searchList []SearchData
        var query = make(map[string]interface{})

        for key, value := range r.URL.Query() {
            switch key {
            case "event":
                event = value[0]
            case "no":
                clientno = value[0]
            case "uid":
                uid = value[0]
            }
        }

        if len(event) > 0 { // event search
            if _, err := strconv.Atoi(event); err == nil {

                query["value."+event] = map[string]interface{}{"$exists":true}
                qs, _ := dbClient.ReadAll(db_name, "lang", query, nil)

                for _, result := range qs {
                    var data = SearchData{}
                    data.DB = "lang"
                    if value, ok := result["value"].(map[string]interface{}); ok {
                        if str, f := value[event].(string); f {data.Value = str}
                    }
                    if str, f := result["key"].(string); f{data.Key = str}
                    
                    searchList = append(searchList, data)
                }
            }
        }

        if len(uid) == 20 { // uid search in both device and mapping
            query["key"] = uid
            uidFunc := func(c string) {
                q := dbClient.ReadOne(db_name, c, query)

                if q != nil {
                    var data = SearchData{}
                    data.DB = c
                    data.Key = uid
                    if value, ok := q["value"].(map[string]interface{}); ok {
                        tmpStr := "{"
                        tmpStr += fmt.Sprint(value)
                        tmpStr += "}"

                        data.Value = tmpStr
                    }
                    
                    searchList = append(searchList, data)
                }
            }
            uidFunc("device")
            uidFunc("mapping")
            uidFunc("allow")
        }

        if len(clientno) > 0 {
            query["key"] = clientno

            q := dbClient.ReadOne(db_name, "client", query)

            if q != nil {
                var data = SearchData{}
                data.DB = "client"
                data.Key = clientno
                if value, ok := q["value"].(map[string]interface{}); ok {
                    tmpStr := "{"
                    tmpStr += fmt.Sprint(value)
                    tmpStr += "}"

                    data.Value = tmpStr
                }
                
                searchList = append(searchList, data)
            }
        }

        input["Data"] = searchList
    } else { // Post
        r.ParseMultipartForm(0)

    }

    t.Execute(w, input)
}
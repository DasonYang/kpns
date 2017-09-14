package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type searchTemplateInfo struct {
	event    string
	clientNo string
	uid      string
}

type searchData struct {
	Key   string
	DB    string
	Value string
}

func (info searchTemplateInfo) getInput() map[string]interface{} {
	var input = make(map[string]interface{})
	var searchList []searchData
	var query = make(map[string]interface{})

	if len(info.event) > 0 { // event search
		if _, err := strconv.Atoi(info.event); err == nil {

			query["value."+info.event] = map[string]interface{}{"$exists": true}
			qs, _ := dbClient.ReadAll(dbName, "lang", query, nil)

			for _, result := range qs {
				var data = searchData{}
				data.DB = "lang"
				if value, ok := result["value"].(map[string]interface{}); ok {
					if str, f := value[info.event].(string); f {
						data.Value = str
					}
				}
				if str, f := result["key"].(string); f {
					data.Key = str
				}

				searchList = append(searchList, data)
			}
		}
	}

	if len(info.uid) == 20 { // uid search in both device and mapping
		query["key"] = info.uid
		uidFunc := func(c string) {
			q := dbClient.ReadOne(dbName, c, query)

			if q != nil {
				var data = searchData{}
				data.DB = c
				data.Key = info.uid
				if value, ok := q["value"].(map[string]interface{}); ok {

					b, err := json.MarshalIndent(value, "", "  ")
					if err != nil {
						fmt.Println("error:", err)
					}
					data.Value = string(b)
				}

				searchList = append(searchList, data)
			}
		}
		uidFunc("device")
		uidFunc("mapping")
		uidFunc("allow")
	}

	if len(info.clientNo) > 0 {
		query["key"] = info.clientNo

		q := dbClient.ReadOne(dbName, "client", query)

		if q != nil {
			var data = searchData{}
			data.DB = "client"
			data.Key = info.clientNo
			if value, ok := q["value"].(map[string]interface{}); ok {
				b, err := json.MarshalIndent(value, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				data.Value = string(b)
			}

			searchList = append(searchList, data)
		}
	}

	input["Data"] = searchList

	return input
}

// SearchHandler - Handle /search
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var input = make(map[string]interface{})

	t, err := template.ParseFiles(templatePath + "/search.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		args := searchTemplateInfo{}

		for key, value := range r.URL.Query() {
			switch key {
			case "event":
				args.event = value[0]
			case "no":
				args.clientNo = value[0]
			case "uid":
				args.uid = value[0]
			}
		}

		input = args.getInput()
	} else { // Post
		r.ParseMultipartForm(0)

	}

	t.Execute(w, input)
}

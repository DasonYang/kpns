package web

import (
	// "log"
	"fmt"
	"time"

	b64 "encoding/base64"
	"html/template"
	"io/ioutil"
	"net/http"

	"kpns/utils"
)

type appKeyTemplateInfo struct {
	displayLimit int
	pageIndex    int
	writable     bool
	query        map[string]interface{}
	popDialog    bool
	dialogMsg    string
}

type appKeyData struct {
	AppID  string
	AppKey string
	Count  int
	Last   string
}

func (info appKeyTemplateInfo) genInput() map[string]interface{} {
	var input = make(map[string]interface{})
	var appList []appKeyData

	qs, count := dbClient.ReadAll(dbName, "appkey", info.query, nil)

	for _, appid := range qs {
		var data appKeyData
		if str, f := appid["key"].(string); f {
			data.AppID = str
		}

		value := appid["value"].(map[string]interface{})

		if ts, f := value["lasttime"].(float64); f {
			tm := time.Unix(int64(ts), 0)
			data.Last = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
		} else if d, f := value["lasttime"].(string); f {
			data.Last = d
		}

		if c, f := value["count"].(int); f {
			data.Count = c
		}

		if k, f := value["appkey"].(string); f {
			data.AppKey = k
		}
		appList = append(appList, data)
	}

	input["Data"] = appList
	input["Count"] = count
	input["Success"] = info.popDialog
	input["Writable"] = info.writable

	return input
}

// AppKeyHandler - Handler /appkey
func AppKeyHandler(w http.ResponseWriter, r *http.Request) {

	// username, password, ok := r.BasicAuth()
	fmt.Println("================================Appkey=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	args := appKeyTemplateInfo{writable: writable, query: make(map[string]interface{})}

	t, err := template.ParseFiles(templatePath + "/appkey.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		fmt.Println("================================Appkey.GET=================================")
		var appID, active string
		for key, value := range r.URL.Query() {
			// fmt.Printf("key = %v, value = %v\n", key, value)
			switch key {
			case "id":
				appID = value[0]
			case "active":
				active = value[0]
			case "key":
				// key = value[0]
			}
		}

		if active == "del" {
			err := dbClient.Delete(dbName, "appkey", map[string]interface{}{"key": appID})
			if err != nil {
				panic(err)
			}
		}

		input := args.genInput()

		t.Execute(w, input)
	} else {
		fmt.Println("================================Allow.POST=================================")
		r.ParseMultipartForm(0)
		var appkey, appID string

		for key, value := range r.PostForm {
			switch key {
			case "id":
				appID = value[0]
			case "key":
				appkey = value[0]
			}
		}

		if _, ok := r.PostForm["search"]; ok { // Search clicked
			var orQuery []map[string]interface{}

			if len(appID) > 0 {
				orQuery = append(orQuery, map[string]interface{}{"key": appID})
			}
			if len(appkey) > 0 {
				orQuery = append(orQuery, map[string]interface{}{"value.appkey": appkey})
			}

			args.query["$or"] = orQuery
		} else if _, ok := r.PostForm["save"]; ok && len(appID) > 0 { // Save clicked
			var data = map[string]interface{}{"key": appID}
			file, _, _ := r.FormFile("File")

			if file != nil { // iOS pem
				defer file.Close()

				dat, err := ioutil.ReadAll(file)

				if err != nil {
					fmt.Printf("Read file with err = %v\n", err)
				}

				sEnc := b64.StdEncoding.EncodeToString(dat)
				strToHash := appkey + "@" + sEnc
				sHash := utils.Hash128(strToHash)

				appkey = appkey + "@" + sHash

				appkeyfile := map[string]interface{}{"key": sHash, "value": sEnc}
				dbClient.Write(dbName, "appkeyfile", appkeyfile, nil)

			}

			if len(appkey) > 0 {
				app := dbClient.ReadOne(dbName, "appkey", map[string]interface{}{"key": appID})
				value := make(map[string]interface{})

				if val, ok := app["value"]; ok { // If data exists
					if _, match := val.(map[string]interface{}); match {
						value = val.(map[string]interface{})
					}
				} else {
					value["first_time"] = int32(time.Now().Unix())
				}

				value["appkey"] = appkey
				data["value"] = value

				dbClient.Write(dbName, "appkey", data, nil)
				args.query["key"] = appID
			}

		}

		input := args.genInput()

		t.Execute(w, input)
	}
}

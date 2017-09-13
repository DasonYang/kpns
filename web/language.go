package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var validLangBatchFile = regexp.MustCompile(`^-?[A-Za-z-_]*\s*(?:,\s*[\d]+\s*(?:,\s*[\w]+\s*)?)?$`)

type LangData struct {
	LangType string
	Status   string
	Message  string
}

func LangHandler(w http.ResponseWriter, r *http.Request) {

	// username, password, ok := r.BasicAuth()
	fmt.Println("================================Lang=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	var query = make(map[string]interface{})
	var displayLang string
	// getInput
	genInput := func(lang string, query map[string]interface{}, success bool, writable bool) map[string]interface{} {
		//fmt.Printf("limit = %v, page = %v, note = %v, query = %v\n", limit, page, note, query)
		var input = make(map[string]interface{})
		var langList []LangData

		// fmt.Printf("params = %v\n", params)
		total := dbClient.Count(db_name, "lang", nil)
		qs, count := dbClient.ReadAll(db_name, "lang", query, nil)

		// fmt.Printf("qs = %v, count = %v\n", qs, count)

		for _, lang := range qs {
			if str, f := lang["key"].(string); f {

				value := lang["value"].(map[string]interface{})

				for k, v := range value {
					data := LangData{}
					data.LangType = str
					data.Status = k
					data.Message = v.(string)

					langList = append(langList, data)
				}
			}

		}

		input["Data"] = langList
		input["Count"] = count
		input["Success"] = success
		input["Writable"] = writable
		input["Total"] = total
		input["Lang"] = lang

		return input
	}
	// genInput
	t, err := template.ParseFiles(TemplatePath + "/lang.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		fmt.Println("================================Lang.GET=================================")
		var active, status string
		lang := "enUS"
		for key, value := range r.URL.Query() {
			fmt.Printf("key = %v, value = %v\n", key, value)
			switch key {
			case "active":
				active = value[0]
			case "status":
				status = value[0]
			case "lang":
				lang = value[0]
			}
		}
		displayLang = lang
		query["key"] = map[string]interface{}{"$regex": lang}

		if active == "del" {
			err := dbClient.Update(db_name,
				"lang",
				map[string]interface{}{"key": lang},
				map[string]interface{}{"$unset": map[string]interface{}{("value." + status): 1}},
				nil)
			if err != nil {
				log.Printf("Delete status %v of %v error : %v\n", status, lang, err)
			}
		}

		input := genInput(displayLang, query, false, writable)

		t.Execute(w, input)
	} else {
		fmt.Println("================================Lang.POST=================================")
		r.ParseMultipartForm(0)

		// var ltime string
		// var mode string
		var isBatchDone bool

		if _, ok := r.PostForm["bsubmit"]; ok { // Handler uploaded file
			file, _, opErr := r.FormFile("bf")

			if opErr != nil {
				fmt.Println(opErr)
			}
			if file != nil {
				defer file.Close()

				dat, readErr := ioutil.ReadAll(file)

				if readErr != nil {
					fmt.Printf("Read file with err = %v\n", readErr)
				}

				lines := strings.Split(string(dat), "\n")

				for _, line := range lines {
					if validLangBatchFile.MatchString(line) {
						cols := strings.Split(line, ",")

						length := len(cols)
						lang := strings.TrimSpace(cols[0])
						var msg, status string

						if length > 1 {
							status = strings.TrimSpace(cols[1])
						}

						if length == 3 {
							msg = strings.TrimSpace(cols[2])
						}

						if ok := strings.HasPrefix(lang, "-"); ok {

							if length == 1 {
								delErr := dbClient.Delete(db_name, "lang", map[string]interface{}{"key": lang[1:]})
								if delErr != nil {
									log.Printf("Error while deleting whole lang : %v, message = %v\n", lang[1:], delErr)
								}
							} else {
								unsetErr := dbClient.Update(db_name,
									"lang",
									map[string]interface{}{"key": lang[1:]},
									map[string]interface{}{"$unset": map[string]interface{}{("value." + status): 1}},
									nil)
								if unsetErr != nil {
									log.Printf("Error while updating data with status : %v, message = %v\n", status, unsetErr)
								}
							}
							continue
						}

						setErr := dbClient.Update(db_name,
							"lang",
							map[string]interface{}{"key": lang},
							map[string]interface{}{"$set": map[string]interface{}{("value." + status): msg}},
							nil)

						if setErr != nil {
							log.Printf("Error while set data with status : %v, message = %v\n", status, setErr)
						}
					}
				}
				isBatchDone = true
			}
		} else { // Handle normal operation
			var lang, status, msg string
			for key, value := range r.PostForm {
				fmt.Printf("key = %v, value = %v\n", key, value)
				switch key {
				case "lang":
					lang = value[0]
				case "status":
					status = value[0]
				case "msg":
					msg = value[0]

				}
			}

			if _, ok := r.PostForm["search"]; ok { // Search clicked

				if len(lang) > 0 {
					query["key"] = map[string]interface{}{"$regex": lang}
				}
			} else if _, ok := r.PostForm["save"]; ok { // Save clicked
				if len(lang) > 0 && len(status) > 0 {
					query["key"] = lang
					setErr := dbClient.Update(db_name,
						"lang",
						map[string]interface{}{"key": lang},
						map[string]interface{}{"$set": map[string]interface{}{("value." + status): msg}},
						nil)

					if setErr != nil {
						log.Printf("Error while set data with status : %v, message = %v\n", status, setErr)
					}
				}
			}
			displayLang = lang
		}
		input := genInput(displayLang, query, isBatchDone, writable)
		t.Execute(w, input)
	}
}

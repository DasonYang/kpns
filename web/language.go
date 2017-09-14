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

type langTemplateInfo struct {
	displayLimit int
	pageIndex    int
	writable     bool
	query        map[string]interface{}
	popDialog    bool
	dialogMsg    string
	lang         string
}

type langData struct {
	LangType string
	Status   string
	Message  string
}

func (info langTemplateInfo) genInput() map[string]interface{} {
	var input = make(map[string]interface{})
	var langList []langData
	total := dbClient.Count(dbName, "lang", nil)
	qs, count := dbClient.ReadAll(dbName, "lang", info.query, nil)

	for _, lang := range qs {
		if str, f := lang["key"].(string); f {

			value := lang["value"].(map[string]interface{})

			for k, v := range value {
				data := langData{}
				data.LangType = str
				data.Status = k
				data.Message = v.(string)

				langList = append(langList, data)
			}
		}
	}

	input["Data"] = langList
	input["Count"] = count
	input["Success"] = info.popDialog
	input["Writable"] = info.writable
	input["Total"] = total
	input["Lang"] = info.lang

	return input
}

// LangHandler -  Handle /lang
func LangHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================================Lang=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	args := langTemplateInfo{writable: writable, query: make(map[string]interface{})}
	t, err := template.ParseFiles(templatePath + "/lang.tmpl")
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
				args.lang = value[0]
			}
		}

		args.query["key"] = map[string]interface{}{"$regex": args.lang}

		if active == "del" {
			err := dbClient.Update(dbName,
				"lang",
				map[string]interface{}{"key": lang},
				map[string]interface{}{"$unset": map[string]interface{}{("value." + status): 1}},
				nil)
			if err != nil {
				log.Printf("Delete status %v of %v error : %v\n", status, lang, err)
			}
		}

		input := args.genInput()

		t.Execute(w, input)
	} else {
		fmt.Println("================================Lang.POST=================================")
		r.ParseMultipartForm(0)

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
								delErr := dbClient.Delete(dbName, "lang", map[string]interface{}{"key": lang[1:]})
								if delErr != nil {
									log.Printf("Error while deleting whole lang : %v, message = %v\n", lang[1:], delErr)
								}
							} else {
								unsetErr := dbClient.Update(dbName,
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

						setErr := dbClient.Update(dbName,
							"lang",
							map[string]interface{}{"key": lang},
							map[string]interface{}{"$set": map[string]interface{}{("value." + status): msg}},
							nil)

						if setErr != nil {
							log.Printf("Error while set data with status : %v, message = %v\n", status, setErr)
						}
					}
				}
				args.popDialog = true
			}
		} else { // Handle normal operation
			var status, msg string
			for key, value := range r.PostForm {
				fmt.Printf("key = %v, value = %v\n", key, value)
				switch key {
				case "lang":
					args.lang = value[0]
				case "status":
					status = value[0]
				case "msg":
					msg = value[0]

				}
			}

			if _, ok := r.PostForm["search"]; ok { // Search clicked

				if len(args.lang) > 0 {
					args.query["key"] = map[string]interface{}{"$regex": args.lang}
				}
			} else if _, ok := r.PostForm["save"]; ok { // Save clicked
				if len(args.lang) > 0 && len(status) > 0 {
					args.query["key"] = args.lang
					setErr := dbClient.Update(dbName,
						"lang",
						map[string]interface{}{"key": args.lang},
						map[string]interface{}{"$set": map[string]interface{}{("value." + status): msg}},
						nil)

					if setErr != nil {
						log.Printf("Error while set data with status : %v, message = %v\n", status, setErr)
					}
				}
			}
		}
		input := args.genInput()
		t.Execute(w, input)
	}
}

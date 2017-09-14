package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validBatchFile = regexp.MustCompile(`^-?[A-Z0-9]{20}\s*,\s*[\d]{4}/[\d]{1,2}/[\d]{1,2}\s*,\s*[\w]+\s*(?:,\s*[\w]+)?\s*$`)

type allowTemplateInfo struct {
	displayLimit int
	pageIndex    int
	writable     bool
	query        map[string]interface{}
	popDialog    bool
	dialogMsg    string
	note         string
	uid          string
}

type allowData struct {
	UID     string
	Updated string
	Limit   string
	Note    string
}

func (info allowTemplateInfo) genInput() map[string]interface{} {
	var input = make(map[string]interface{})
	var params = make(map[string]interface{})
	var allowList []allowData
	pageIdx := info.pageIndex
	displayLimit := info.displayLimit
	query := info.query

	if displayLimit < 20 {
		displayLimit = 20
	}
	if pageIdx == 0 {
		pageIdx = 1
	}

	params["skip"] = (pageIdx - 1) * displayLimit
	params["limit"] = displayLimit

	// fmt.Printf("params = %v\n", params)

	qs, count := dbClient.ReadAll(dbName, "allow", query, params)

	for _, allow := range qs {
		var data allowData
		if str, f := allow["key"].(string); f {
			data.UID = str
		}

		value := allow["value"].(map[string]interface{})

		if str, f := value["update_time"].(string); f {
			data.Updated = str
		}

		if ts, f := value["limit"].(float64); f {
			tm := time.Unix(int64(ts), 0)
			data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
		} else if ts, f := value["limit"].(int); f {
			tm := time.Unix(int64(ts), 0)
			data.Limit = fmt.Sprintf("%v", tm.Format("2006-01-02 15:04:05"))
		}

		if str, f := value["note"].(string); f {
			data.Note = str
		}

		// fmt.Printf("type of limit = %v\n", reflect.TypeOf(value["limit"]))
		allowList = append(allowList, data)
	}

	// fmt.Printf("count = %v\n", count)

	input["Data"] = allowList
	input["Page"] = pageIdx
	input["Count"] = count
	input["Limit"] = displayLimit
	input["HasNote"] = true
	input["Note"] = info.note
	input["Success"] = info.popDialog
	input["Writable"] = info.writable

	if pageIdx > 1 {
		input["HasPre"] = true
		input["Pre"] = pageIdx - 1
	}
	if (pageIdx * displayLimit) < count {
		input["HasNext"] = true
		input["Next"] = pageIdx + 1
	}

	return input
}

// AllowHandler - Handler /allow
func AllowHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("================================Allow=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	args := allowTemplateInfo{writable: writable, query: make(map[string]interface{})}

	t, err := template.ParseFiles(templatePath + "/allow.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		fmt.Println("================================Allow.GET=================================")
		var active string
		for key, value := range r.URL.Query() {
			fmt.Printf("key = %v, value = %v\n", key, value)
			switch key {
			case "page":
				args.pageIndex, _ = strconv.Atoi(value[0])
			case "limit":
				args.displayLimit, _ = strconv.Atoi(value[0])
			case "active":
				active = value[0]
			case "uid":
				args.uid = value[0]
			case "note":
				args.note = value[0]
				args.query["value.note"] = map[string]interface{}{"$regex": args.note}
			}
		}

		if active == "del" && len(args.uid) == 20 {
			err := dbClient.Delete(dbName, "allow", map[string]interface{}{"key": args.uid})
			if err != nil {
				panic(err)
			}
		}

		input := args.genInput()

		t.Execute(w, input)
	} else {
		fmt.Println("================================Allow.POST=================================")
		r.ParseMultipartForm(0)

		var ltime string
		var mode string

		if _, ok := r.PostForm["bsubmit"]; ok { // Handler uploaded file
			file, _, err := r.FormFile("bf")

			if err != nil {
				fmt.Println(err)
			}
			if file != nil {
				defer file.Close()

				dat, err := ioutil.ReadAll(file)

				if err != nil {
					fmt.Printf("Read file with err = %v\n", err)
				}

				lines := strings.Split(string(dat), "\n")

				for _, line := range lines {
					if validBatchFile.MatchString(line) {
						cols := strings.Split(line, ",")
						_uid := strings.TrimSpace(cols[0])

						if ok := strings.HasPrefix(_uid, "-"); ok {
							err := dbClient.Delete(dbName, "allow", map[string]interface{}{"key": _uid[1:]})
							if err != nil {
								log.Printf("Error while deleting data with uid : %v, message = %v\n", _uid[1:], err)
							}
							continue
						}

						info := map[string]interface{}{"key": _uid}
						data := make(map[string]interface{})

						data["note"] = strings.TrimSpace(cols[2])
						tm, _ := time.Parse("2006/01/02", strings.TrimSpace(cols[1]))
						data["limit"] = int32(tm.Unix())
						if len(cols) == 4 {
							data["method"] = strings.TrimSpace(cols[3])
						}
						data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
						info["value"] = data

						err = dbClient.Write(dbName, "allow", info, nil)

						if err != nil {
							log.Printf("Error while inserting data with uid : %v, message = %v\n", _uid[1:], err)
						}
					}
				}
				args.popDialog = true
			}
		} else { // Handle normal operation
			for key, value := range r.PostForm {
				fmt.Printf("key = %v, value = %v\n", key, value)
				switch key {
				case "ltime":
					ltime = value[0]
				case "limit":
					args.displayLimit, _ = strconv.Atoi(value[0])
				case "note":
					args.note = value[0]
					args.query["value.note"] = map[string]interface{}{"$regex": args.note}
				case "File":
				case "mode":
					mode = value[0]
				case "uid":
					args.uid = value[0]

				}
			}

			if _, ok := r.PostForm["search"]; ok { // Search clicked
				if len(args.uid) > 0 {
					args.query["key"] = args.uid
				}
			} else if _, ok := r.PostForm["save"]; ok { // Save clicked

				if len(args.uid) > 0 && len(ltime) > 0 && len(args.note) > 0 {
					tm, err := time.Parse("2006/01/02", ltime)

					if err != nil {
						panic(err)
					}
					fmt.Println(tm.UnixNano(), ltime, mode)

					args.query["key"] = args.uid

					info := map[string]interface{}{"key": args.uid}
					data := make(map[string]interface{})
					data["limit"] = int32(tm.Unix())
					data["note"] = args.note
					data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
					info["value"] = data
					dbClient.Write(dbName, "allow", info, nil)
					args.note = ""
				} else {
					args.note = ""
				}
			}
		}

		input := args.genInput()

		t.Execute(w, input)
	}
}

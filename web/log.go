package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type logData struct {
	DateTime string
	Status   string
	Content  string
}

type logTemplateInfo struct {
	displayLimit int
	pageIndex    int
	writable     bool
	query        map[string]interface{}
	popDialog    bool
	dialogMsg    string
	from         string
	to           string
	ip           string
	status       string
	text         string
}

func (info logTemplateInfo) genInput() map[string]interface{} {
	var input = make(map[string]interface{})
	var params = make(map[string]interface{})
	var logList []logData

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
	params["sort"] = "-$natural"

	qs, count := dbClient.ReadAll(dbName, "logs", query, params)

	for _, log := range qs {
		var data logData
		value := log["value"].(map[string]interface{})

		if str, f := value["date"].(string); f {
			var processTime float64
			if tm, g := value["time"].(float64); g {
				processTime = tm
			}
			data.DateTime = fmt.Sprintf("%v(%.2f)", str, processTime)
		} else if t, f := value["date"].(time.Time); f {
			var processTime float64
			if tm, g := value["time"].(float64); g {
				processTime = tm
			}
			data.DateTime = fmt.Sprintf("%v(%.2f)", t.Format("2006-01-02 15:04:05"), processTime)
		}

		if str, f := value["status"].(string); f {
			var hostIP = "__"
			if h, g := value["host"]; g {
				hostIP = fmt.Sprintf("%v", h)
			}
			data.Status = fmt.Sprintf("%v > %v", str, hostIP)
		}

		if msg, f := value["msg"].(string); f {
			var address = "NoIP"
			var ct string
			if addr, ook := value["address"].(string); ook {
				address = addr
			}

			if content, ook := value["data"]; ook {
				b, err := json.MarshalIndent(content, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}

				ct = string(b)
			}
			data.Content = fmt.Sprintf("%v @ %v >> %v", msg, address, ct)
		}

		logList = append(logList, data)
	}

	input["Data"] = logList
	input["Page"] = pageIdx
	input["Count"] = count
	input["Limit"] = displayLimit
	input["Success"] = info.popDialog
	input["Writable"] = info.writable
	input["Msg"] = info.dialogMsg
	input["To"] = info.to
	input["Ip"] = info.ip
	input["Status"] = info.status
	input["Text"] = info.text

	if len(info.from) > 0 {
		input["From"] = info.from
	} else {
		yesterday := time.Now().AddDate(0, 0, -1)
		input["From"] = yesterday.Format("2006/01/02 15:04:05")
	}

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

// LogHandler - Log handler
func LogHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("================================Account=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	// var query = make(map[string]interface{})
	args := logTemplateInfo{writable: writable, query: make(map[string]interface{})}
	// genInput
	t, err := template.ParseFiles(templatePath + "/log.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		fmt.Println("================================Log.GET=================================")
		var keyRange = make(map[string]interface{})
		layout := "2006/01/02 15:04:05"

		for key, value := range r.URL.Query() {
			fmt.Printf("key = %v, value = %v\n", key, value)
			switch key {
			case "page":
				args.pageIndex, _ = strconv.Atoi(value[0])
			case "limit":
				args.displayLimit, _ = strconv.Atoi(value[0])
			case "from":
				fromStr := value[0]

				if len(fromStr) > 0 {
					fromDate, err := time.Parse(layout, fromStr)
					fmt.Printf("err = %v\n", err)

					keyRange["$gte"] = fromDate
					fmt.Printf("fromDate = %v\n", fromDate)
				}
				args.from = fromStr
			case "to":
				toStr := value[0]
				if len(toStr) > 0 {
					toDate, err := time.Parse(layout, toStr)
					fmt.Printf("err = %v\n", err)

					keyRange["$lte"] = toDate
					fmt.Printf("toDate = %v\n", toDate)
				}
				args.to = toStr
			case "text":
				txtStr := value[0]
				if len(txtStr) > 0 {

				}
				args.text = txtStr
			case "status":
				status := value[0]

				if len(status) > 0 {
					args.query["value.status"] = map[string]interface{}{"$regex": status}
				}

				args.status = status
			case "ip":
				ipStr := value[0]

				if len(ipStr) > 0 {
					args.query["value.address"] = map[string]interface{}{"$regex": ipStr}
				}
				args.ip = ipStr
			}
		}

		if len(keyRange) > 0 {
			args.query["key"] = keyRange
		}

		input := args.genInput()

		t.Execute(w, input)
	} else {
		fmt.Println("================================Log.POST=================================")
		r.ParseMultipartForm(0)

		t.Execute(w, nil)
	}
}

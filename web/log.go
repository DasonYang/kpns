package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// LogData : Log data format
type LogData struct {
	DateTime string
	Status   string
	Content  string
}

// LogHandler - Log handler
func LogHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("================================Account=================================", r.Context().Value("Writable"))
	var writable = r.Context().Value("Writable").(bool)
	var query = make(map[string]interface{})
	var pageIdx, limit int
	// // var account string
	var resultMsg string
	var popDialog bool
	// getInput
	genInput := func(args map[string]interface{}) map[string]interface{} {
		var input = make(map[string]interface{})
		var params = make(map[string]interface{})
		var logList []LogData
		pageIdx := args["page"].(int)
		displayLimit := args["limit"].(int)
		query := args["query"].(map[string]interface{})

		if limit < 20 {
			displayLimit = 20
		}
		if pageIdx == 0 {
			pageIdx = 1
		}

		params["skip"] = (pageIdx - 1) * displayLimit
		params["limit"] = displayLimit
		params["sort"] = "-$natural"

		qs, count := dbClient.ReadAll(db_name, "logs", query, params)

		for _, log := range qs {
			var data LogData
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

				if addr, g := value["address"].(string); g {
					address = addr
				}
				data.Content = fmt.Sprintf("%v @ %v >> %v", msg, address, value["data"])
				data.Content = strings.Replace(data.Content, "map", "", -1)
			}

			logList = append(logList, data)
		}

		input["Data"] = logList
		input["Page"] = pageIdx
		input["Count"] = count
		input["Limit"] = displayLimit
		input["Success"] = args["needPop"].(bool)
		input["Writable"] = args["writable"].(bool)
		input["Msg"] = args["popMsg"].(string)
		input["To"] = args["to"].(string)
		input["Ip"] = args["ip"].(string)
		input["Status"] = args["status"].(string)
		input["Text"] = args["text"].(string)

		fromStr := args["from"].(string)
		if len(fromStr) > 0 {
			input["From"] = fromStr
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
	// genInput
	t, err := template.ParseFiles(TemplatePath + "/log.tmpl")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		panic(err)
	}

	if r.Method == "GET" {
		fmt.Println("================================Log.GET=================================")
		var keyRange = make(map[string]interface{})
		var args = map[string]interface{}{"page": 0, "limit": 0, "from": "", "to": "", "ip": "", "text": "", "status": ""}
		layout := "2006/01/02 15:04:05"
		for key, value := range r.URL.Query() {
			fmt.Printf("key = %v, value = %v\n", key, value)
			switch key {
			case "page":
				pageIdx, _ = strconv.Atoi(value[0])
			case "limit":
				limit, _ = strconv.Atoi(value[0])
			case "from":
				fromStr := value[0]

				if len(fromStr) > 0 {
					fromDate, err := time.Parse(layout, fromStr)
					fmt.Printf("err = %v\n", err)

					keyRange["$gte"] = fromDate
					fmt.Printf("fromDate = %v\n", fromDate)
				}
				args["from"] = fromStr
			case "to":
				toStr := value[0]
				if len(toStr) > 0 {
					toDate, err := time.Parse(layout, toStr)
					fmt.Printf("err = %v\n", err)

					keyRange["$lte"] = toDate
					fmt.Printf("toDate = %v\n", toDate)
				}
				args["to"] = toStr
			case "text":
				txtStr := value[0]
				if len(txtStr) > 0 {

				}
				args["text"] = txtStr
			case "status":
				status := value[0]

				if len(status) > 0 {
					query["value.status"] = map[string]interface{}{"$regex": status}
				}

				args["status"] = status
			case "ip":
				ipStr := value[0]

				if len(ipStr) > 0 {
					query["value.address"] = map[string]interface{}{"$regex": ipStr}
				}
				args["ip"] = ipStr
			}
		}

		if len(keyRange) > 0 {
			query["key"] = keyRange
		}

		args["page"] = pageIdx
		args["limit"] = limit
		args["query"] = query
		args["needPop"] = popDialog
		args["popMsg"] = resultMsg
		args["writable"] = writable

		input := genInput(args)

		t.Execute(w, input)
	} else {
		fmt.Println("================================Log.POST=================================")
		r.ParseMultipartForm(0)

		t.Execute(w, nil)
	}
}

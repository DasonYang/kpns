package web

import(
    "fmt"

    "net/http"
    "html/template"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles(TemplatePath+"/search.tmpl")
    if err != nil {
        fmt.Printf("Error = %v\n", err)
        panic(err)
    }

    t.Execute(w, nil)
}
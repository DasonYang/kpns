package kpns

import (
    "fmt"
    "net/http"
    "html/template"

)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintln(w, "200 Success"+" "+r.URL.Path[1:])
    if r.Method == "GET" {
        // info := make(map[string]interface{})
        t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
        if err != nil {
            fmt.Printf("Error = %v\n", err)
            panic(err)
        }
        t.Execute(w, nil)
    } else {
        r.ParseMultipartForm(0)
        // logic part of log in
        username := r.PostForm["user"][0]
        password := r.PostForm["pswd"][0]
        if len(username) == 0 || len(password) == 0 {
            fmt.Println("username is empty")
            info := make(map[string]interface{})
            info["Error"] = true
            info["Message"] = "Login Error!"

            t, err := template.ParseFiles(TemplatePath+"/login.tmpl")
            if err != nil {
                fmt.Printf("Error = %v\n", err)
                panic(err)
            }
            t.Execute(w, info)
        } else {
            http.Redirect(w, r, "/allow", http.StatusSeeOther)
        }
    }
}
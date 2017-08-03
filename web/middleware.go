package web

import(
    // "fmt"
    "net/http"
)

var (
    Permissions = map[string][]string{"all":{"allow", "appkey", "search"}, 
                                      "editor":{"search",},
                                      "basic":{"search",}}
)

type Adapter func(http.Handler) http.Handler

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var user, token string

        var cookies = make(map[string]string)
        for _, cookie := range r.Cookies() {cookies[cookie.Name] = cookie.Value}

        if val, ok := cookies["user"]; ok {user = val} 
        if val, ok := cookies["token"]; ok {token = val}

        if user == "" || token == "" {
            // Not login or expired
            http.Redirect(w, r, "/login", http.StatusSeeOther)
        } else {
            // Auth verified
            ret := dbClient.ReadOne("tpns", "account", map[string]interface{}{"key":user})
            if val, ok := ret["value"]; ok {
                value := val.(map[string]interface{})

                if t, ok := value["token"].(string); ok {
                    if cookies["token"] != t {
                        http.SetCookie(w, &http.Cookie{Name: "token", Value: ""})
                        http.SetCookie(w, &http.Cookie{Name: "mode", Value: ""})
                        http.SetCookie(w, &http.Cookie{Name: "user", Value: ""})
                        http.SetCookie(w, &http.Cookie{Name: "msg", Value: "You already login in other device."})
                        http.Redirect(w, r, "/login", http.StatusSeeOther)
                    }
                }

            }
        }

        next.ServeHTTP(w, r)
    })
}

func Middlewares(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        for _, h := range middlewares {
            next = h(next)
        }
    })
}

func Adapt(next http.Handler, middlewares ...Adapter) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        for _, h := range middlewares {
            next = h(next)
        }
    })
}
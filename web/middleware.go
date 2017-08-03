package web

import(
    "fmt"
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
        fmt.Println("Executing kpnsMiddleware")
        var user, token string
        path := r.URL.Path[1:]

        var cookies = make(map[string]string)
        for _, cookie := range r.Cookies() {
            fmt.Printf("Cookie : %v\n", cookie)
            cookies[cookie.Name] = cookie.Value

        }
        fmt.Printf("Cookie : %v\n", cookies)

        if val, ok := cookies["user"]; ok {
            user = val
        } 

        if val, ok := cookies["token"]; ok {
            token = val
        }

        if user == "" || token == "" {
            // Not login or expired
            if path != "login" {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
            }
        } else {
            // Auth verified
            ret := dbClient.ReadOne("tpns", "account", map[string]interface{}{"key":user})
            if val, ok := ret["value"]; ok {
                value := val.(map[string]interface{})
                fmt.Printf("value = %v\n", value)

                if t, ok := value["token"].(string); ok {
                    if cookies["token"] != t {
                        
                        // expiration := time.Now()
                        // expiration = expiration.Add(time.Minute * time.Duration(1))
                        // fmt.Printf("expiration = %v, token = %v, mode = %v\n", expiration, token, mode)
                        // // cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
                        // http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: expiration})
                        // http.SetCookie(w, &http.Cookie{Name: "mode", Value: mode, Expires: expiration})
                        // http.SetCookie(w, &http.Cookie{Name: "user", Value: username, Expires: expiration})
                        // http.SetCookie(w, &http.Cookie{Name: "msg", Value: "You already login in other device."})
                        // http.Redirect(w, r, "/login", http.StatusSeeOther)
                    }
                }

            }
            if path == "login" {
                http.Redirect(w, r, "/allow", http.StatusSeeOther)
            }
        }

        next.ServeHTTP(w, r)
        fmt.Println("Executing kpnsMiddleware again")
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
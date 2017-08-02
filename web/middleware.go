package web

import(
    "fmt"
    "net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Executing kpnsMiddleware\n")
        var cookies = make(map[string]string)
        for _, cookie := range r.Cookies() {
            fmt.Printf("Cookie : %v\n", cookie)
            cookies[cookie.Name] = cookie.Value

        }
        fmt.Printf("Cookie : %v\n", cookies)
        next.ServeHTTP(w, r)
        fmt.Println("Executing kpnsMiddleware again")
      })
}

// func Middlewares(middlewares ...func(http.Handler) http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         for _, h := range middlewares {
            
//         }
//     })
// }
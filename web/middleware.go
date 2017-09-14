package web

import (
	"context"
	"fmt"
	"net/http"
)

type Adapter func(http.Handler) http.Handler

// AuthMiddleware Middleware for auth
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user, token, mode string
		var allowed bool
		var ctx = context.WithValue(r.Context(), "Writable", false)

		var cookies = make(map[string]string)
		for _, cookie := range r.Cookies() {
			cookies[cookie.Name] = cookie.Value
		}

		if val, ok := cookies["user"]; ok {
			user = val
		}
		if val, ok := cookies["token"]; ok {
			token = val
		}

		if user == "" || token == "" {
			// Not login or expired
			redirect := "/login?next=" + r.URL.Path[1:]
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return
		}
		// Auth verified
		ret := dbClient.ReadOne(dbName, "account", map[string]interface{}{"key": user})
		if val, ok := ret["value"]; ok {
			value := val.(map[string]interface{})

			if t, ok := value["token"].(string); ok {
				if cookies["token"] != t {
					http.SetCookie(w, &http.Cookie{Name: "token", Value: ""})
					http.SetCookie(w, &http.Cookie{Name: "mode", Value: ""})
					http.SetCookie(w, &http.Cookie{Name: "user", Value: ""})
					http.SetCookie(w, &http.Cookie{Name: "msg", Value: "You already login in other device."})
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
			}
			if m, ok := value["mode"].(string); ok {
				mode = m
			}
		}

		if levels, ok := permissions[mode]; ok {

			fmt.Println("val = ", levels)

			if rules, ok := levels[r.URL.Path[1:]]; ok {
				if rules.Readable {
					fmt.Println("allow to access")
					allowed = true

					ctx = context.WithValue(r.Context(), "Writable", rules.Writable)
				}
			}
		}

		if !allowed {
			http.Redirect(w, r, "/search", http.StatusSeeOther)
			return
		}

		fmt.Println("Function : AuthMiddleware, Execute ServeHTTP")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func Middlewares(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		for _, h := range middlewares {
// 			next = h(next)
// 		}
// 	})
// }

// func Adapt(next http.Handler, middlewares ...Adapter) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		for _, h := range middlewares {
// 			next = h(next)
// 		}
// 	})
// }

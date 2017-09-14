package web

import (
	"fmt"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	user, err := r.Cookie("user")

	fmt.Printf("user = %v\n", user)

	if err == nil {
		userData := dbClient.ReadOne(dbName, "account", map[string]interface{}{"key": user.Value})

		if value, ok := userData["value"]; ok {

			value.(map[string]interface{})["token"] = ""
			userData["value"] = value
			dbClient.Write(dbName, "account", userData, nil)
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

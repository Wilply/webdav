package main

import (
	"net/http"
)

func printcookie(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, c := range cookies {
		logger(0, "cookie name : ", c.Name, " ; cookie value : ", c.Value)
		w.Write([]byte("cookie name : " + c.Name + " ; cookie value : " + c.Value))
	}
}

func setcookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "testcookie",
		Value: "testvalue",
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Cookie Set"))
}

package main

import (
	"html/template"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/homepage.html", "./template/head.html", "./template/stylesheet.html")
	if err != nil {
		logger(3, "Cannot parse template for homapage \n", err.Error())
	}

	data := struct {
		Name string
	}{"John Smith"}

	err = t.Execute(w, data)
	if err != nil {
		logger(3, "Cannot execute template for homapage \n", err.Error())
	}
}

func listuser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Login string
		Rw    string
		Ro    string
	}

	var userlist []User

	_, loginlist := getuserlist()

	for _, u := range loginlist {
		_, name, _, rw, ro := getuser(u)
		user := User{
			Login: name,
			Rw:    rw,
			Ro:    ro,
		}
		userlist = append(userlist, user)
	}

	t, err := template.ParseFiles("./template/userlist.html", "./template/head.html", "./template/stylesheet.html")
	if err != nil {
		logger(3, "Cannot parse template for listuser \n", err.Error())
	}

	err = t.Execute(w, userlist)
	if err != nil {
		logger(3, "Cannot execute template for listuser \n", err.Error())
	}
}

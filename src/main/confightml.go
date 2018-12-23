package main

import (
	"html/template"
	"net/http"
)

func webtemplate() (t *template.Template) {
	t, err := template.ParseFiles("./template/updateuser.html", "./template/deleteuser.html", "./template/success.html", "./template/createuser.html", "./template/error.html", "./template/homepage.html", "./template/userlist.html", "./template/head.html", "./template/stylesheet.html")
	if err != nil {
		logger(3, "Cannot parse template file \n", err.Error())
	}
	return
}

func webhomepage(w http.ResponseWriter, r *http.Request) {

	foo := struct { //temporaire
		Name string
	}{"bar"}

	err := webtemplate().ExecuteTemplate(w, "homepage", foo)
	if err != nil {
		logger(3, "Cannot execute \"homepage\" template \n", err.Error())
	}
}

func weblistuser(w http.ResponseWriter, r *http.Request) {

	var userlist []string

	_, loginlist := getuserlist()

	for _, u := range loginlist {
		_, name, _ := getuser(u)
		userlist = append(userlist, name)
	}

	err := webtemplate().ExecuteTemplate(w, "userlist", userlist)
	if err != nil {
		logger(3, "Cannot execute \"listuser\" template \n", err.Error())
	}
}

func webcreateuser(w http.ResponseWriter, r *http.Request) {
	foo := struct {
		Name string
	}{"bar"}

	if r.Method == "GET" { //send the form
		err := webtemplate().ExecuteTemplate(w, "createuser", foo)
		if err != nil {
			logger(3, "Cannot execute \"createuser\" template \n", err.Error())
		}
	} else if r.Method == "POST" { //process the form
		u := r.FormValue("login")
		p := r.FormValue("password")
		if u != "" && p != "" {
			ok := insertuser(u, hashpassword(p))
			if !ok {
				weberror(w, r)
			} else {
				websucces(w, r)
			}
		} else {
			weberror(w, r)
		}
	} else {
		weberror(w, r)
	}
}

func webdeleteuser(w http.ResponseWriter, r *http.Request) {
	var userlist []string

	_, loginlist := getuserlist()

	for _, u := range loginlist {
		if u != "admin" { //DONT MESS WITH ADMIN
			_, name, _ := getuser(u)
			userlist = append(userlist, name)
		}
	}

	if r.Method == "GET" { //send the form
		err := webtemplate().ExecuteTemplate(w, "deleteuser", userlist)
		if err != nil {
			logger(3, "Cannot execute \"deleteuser\" template \n", err.Error())
		}
	} else if r.Method == "POST" { //process
		username := r.FormValue("login")
		ok := deleteuser(username)
		if !ok {
			weberror(w, r)
		} else {
			websucces(w, r)
		}
	} else {
		weberror(w, r)
	}
}

func webupdateuser(w http.ResponseWriter, r *http.Request) {
	var userlist []string

	_, loginlist := getuserlist()

	for _, u := range loginlist {
		_, name, _ := getuser(u)
		userlist = append(userlist, name)
	}

	if r.Method == "GET" { //send the form
		err := webtemplate().ExecuteTemplate(w, "updateuser", userlist)
		if err != nil {
			logger(3, "Cannot execute \"updateuser\" template \n", err.Error())
		}
	} else if r.Method == "POST" {
		u := r.FormValue("login")
		p := r.FormValue("password")
		if u != "" && p != "" {
			ok := updateuserpassword(u, hashpassword(p))
			if !ok {
				weberror(w, r)
			} else {
				websucces(w, r)
			}
		} else {
			weberror(w, r)
		}
	} else {
		weberror(w, r)
	}
}

func weberror(w http.ResponseWriter, r *http.Request) {
	foo := struct { //temporaire
		Name string
	}{"bar"}

	err := webtemplate().ExecuteTemplate(w, "error", foo)
	if err != nil {
		logger(3, "Cannot execute \"error\" template \n", err.Error())
	}
}

func websucces(w http.ResponseWriter, r *http.Request) {
	foo := struct {
		URL string
	}{r.URL.RawPath}

	err := webtemplate().ExecuteTemplate(w, "success", foo)
	if err != nil {
		logger(3, "Cannot execute \"success\" template \n", err.Error())
	}
}

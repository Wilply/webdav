package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/webdav"
)

func createdavhandler() (davhandler *webdav.Handler) {
	return &webdav.Handler{ //on creer un pointer pour mes acces en temps reelle et le locker
		Prefix:     config.DavPrefix,
		FileSystem: webdav.Dir(config.DavRoot),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) { //ecrire les erreurs dans la console
			if err != nil {
				fmt.Printf("%-8s %s %s %s \n", "[ERROR]", r.Method, r.URL, err)
				logger(3, r.Method, r.URL.EscapedPath(), err.Error())
			} else {
				//fmt.Printf("%-8s %s %s \n", "[INFO]", r.Method, r.URL)
				logger(1, r.Method, r.URL.EscapedPath())
			}
		},
	}
}

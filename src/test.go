package test

import (
	"fmt"
	"net/http"
)

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("toto"))
	//http.ServeFile(w, r, "test.txt")
}

var serverhttp http.Handler = http.FileServer(http.Dir("static/"))

func main() {
	fmt.Println("coucou")
	http.Handle("/test/", http.StripPrefix("/test/", serverhttp))
	//http.Handle("/", serverhttp)
	http.HandleFunc("/", test)
	http.ListenAndServe(":8080", nil)
}

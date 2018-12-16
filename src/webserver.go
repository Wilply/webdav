package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/net/webdav"
	"gopkg.in/yaml.v2"
)

type Configs struct {
	LogLevel   int    `yaml:"LogLevel"`
	ServerAddr string `yaml:"ServerAddr"`
	ServerPort int    `yaml:"ServerPort"`
	DavRoot    string `yaml:"DavRoot"`   //chemin local
	DavPrefix  string `yaml:"DavPrefix"` //prefix url
}

var defaultconf = Configs{
	LogLevel:   2,
	ServerAddr: "",
	ServerPort: 8080,
	DavRoot:    "./static/",
	DavPrefix:  "/dav/", // '/' IS NEDEED BEFORE AND AFTER PATH
}

//if not config file provided, fall back to default config
//you don't have to provide all values in config files

var username string = "toto"
var password string = "tata"

var config Configs = defaultconf //copy default config

var serverlistenon string //base url for the server

var davhandler http.Handler //webdav handler

func init() {
	loginfo("Reading & Parsing Config")
	var ok bool
	ok = parseyaml()
	if !ok {
		fmt.Printf("%-8s %s \n", "[ERROR]", "Falling back to default config")
		config = defaultconf
	}

	printconfig()

	serverlistenon = config.ServerAddr + ":" + strconv.Itoa(config.ServerPort)

	davhandler = &webdav.Handler{ //on creer un pointer pour mes acces en temps reelle et le locker
		Prefix:     config.DavPrefix,
		FileSystem: webdav.Dir(config.DavRoot),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) { //ecrire les erreurs dans la console
			if err != nil {
				fmt.Printf("%-8s %s %s %s \n", "[ERROR]", r.Method, r.URL, err)
			} else {
				loginfo(r.Method, r.URL.EscapedPath())
			}
		},
	}

	loginfo("Registrering Handlers") //on ajoute les handlers
	//http.Handle(config.DavPrefix, davhandler) //test without auth
	http.Handle(config.DavPrefix, authandler(davhandler))
}

func main() {
	loginfo("Starting Server")
	fmt.Printf("%-8s %s %s \n", "[INFO]", "Listening on", serverlistenon)
	fmt.Println("*********************************")
	http.ListenAndServe(serverlistenon, nil) //on demarre le serveur
}

func authandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //handlerFunc AVEC UN R, s'occupe de passer les arguments a la func
		//ici on fait des truc => c'est le middleware
		//var username, password, ok = r.BasicAuth()
		var user, pass, ok = r.BasicAuth()
		if !ok { // header authorization no set
			loginfo("connection attempt from ", r.RemoteAddr)
			w.Header().Set("WWW-Authenticate", `Basic realm="webdav acces, auth required" charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			if authenticate(user, pass) { //fonction d'authentification
				loginfo("User : ", user, " is connected from ", r.RemoteAddr)
				//w.WriteHeader(http.StatusOK)
			} else {
				loginfo("Failed connexion attempt from ", r.RemoteAddr, " with user : ", user)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func authenticate(user, pass string) bool { //on modifiera ca plus tard
	if user == username && pass == password {
		return true
	}
	return false
}

func loginfo(msg ...string) { //affiche les info si config.LogLevel > 1
	if config.LogLevel > 1 {
		var str string
		for _, v := range msg {
			str = str + v
		}
		fmt.Printf("%-8s %s \n", "[INFO]", str)
	}
}

func parseyaml() bool {
	var ok = true

	source, readerr := ioutil.ReadFile("config.yaml")

	if readerr != nil {
		ok = false
		fmt.Printf("%-8s %s \n", "[ERROR]", readerr)
		fmt.Printf("%-8s %s \n", "[ERROR]", "Error while reading config file")
		return ok
	}

	parseerr := yaml.UnmarshalStrict(source, &config) //existing key override default
	if parseerr != nil {
		ok = false
		fmt.Printf("%-8s %s \n", "[ERROR]", parseerr)
		fmt.Printf("%-8s %s \n", "[ERROR]", "Error while parsing config file")
		return ok
	}

	return ok
}

func printconfig() {
	if config.LogLevel > 1 {
		fmt.Println("Running-Config : ", config)
	}
}

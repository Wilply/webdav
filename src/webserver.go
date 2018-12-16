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

/* LOG LVL :
CRITICAL = 4
WARNING = 3
INFO = 2
DEBUG = 1
DEBUG++ = 0
*/

//if not config file provided, fall back to default config
//you don't have to provide all values in config files

var username string = "toto"
var password string = "tata"

var config Configs = defaultconf //copy default config

var serverlistenon string //base url for the server

var davhandler http.Handler //webdav handler

func init() {
	consoleapp(1, "Reading & Parsing Config")
	var ok bool
	ok = parseyaml()
	if !ok {
		consoleapp(3, "Falling back to default config")
		config = defaultconf
	}

	printconfig() //config.LogLevel 0

	serverlistenon = config.ServerAddr + ":" + strconv.Itoa(config.ServerPort)

	davhandler = &webdav.Handler{ //on creer un pointer pour mes acces en temps reelle et le locker
		Prefix:     config.DavPrefix,
		FileSystem: webdav.Dir(config.DavRoot),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) { //ecrire les erreurs dans la console
			if err != nil {
				fmt.Printf("%-8s %s %s %s \n", "[ERROR]", r.Method, r.URL, err)
				consoleapp(3, r.Method, r.URL.EscapedPath(), err.Error())
			} else {
				//fmt.Printf("%-8s %s %s \n", "[INFO]", r.Method, r.URL)
				consoleapp(1, r.Method, r.URL.EscapedPath())
			}
		},
	}

	consoleapp(1, "Registrering Handlers") //on ajoute les handlers
	//http.Handle(config.DavPrefix, davhandler) //test without auth
	http.Handle(config.DavPrefix, authandler(davhandler))
	http.HandleFunc("/cookielist", printcookie)
	http.HandleFunc("/cookieset", setcookie)
}

func main() {
	consoleapp(-1, "Starting Server")
	consoleapp(-1, "Listening on", serverlistenon)
	fmt.Println("*********************************")
	http.ListenAndServe(serverlistenon, nil) //on demarre le serveur
}

func authandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //handlerFunc AVEC UN R, s'occupe de passer les arguments a la func
		//1st test if cookie
		c, err := r.Cookie("AuthToken")
		if err == nil {
			consoleapp(0, "Auth cookie ok, auth successful. Token = "+c.Value)
			next.ServeHTTP(w, r)
			return
		} else {
			consoleapp(1, "No cookie for "+r.RemoteAddr+" try to authenticate")
		}
		//ici on fait des truc => c'est le middleware
		//var username, password, ok = r.BasicAuth()
		var user, pass, ok = r.BasicAuth()
		if !ok { // header authorization no set
			consoleapp(2, "connection attempt from ", r.RemoteAddr)
			w.Header().Set("WWW-Authenticate", `Basic realm="webdav acces, auth required" charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			if authenticate(user, pass) { //fonction d'authentification
				authc := http.Cookie{
					Name:  "AuthToken",
					Value: "ziefbi954fze",
				}
				http.SetCookie(w, &authc)
				consoleapp(2, "User : ", user, " is connected from ", r.RemoteAddr)
				//w.WriteHeader(http.StatusOK)
			} else {
				consoleapp(3, "Failed connexion attempt from ", r.RemoteAddr, " with user : ", user)
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

func consoleapp(loglvl int, msg ...string) { //affiche les info si config.LogLevel > 1
	var prefix string
	if loglvl == -1 {
		loglvl = 100
	}

	switch loglvl {
	case 4:
		prefix = "[CRITICAL]"
	case 3:
		prefix = "[WARNING]"
	case 2:
		prefix = "[INFO]"
	case 1:
		prefix = "[DEBUG]"
	case 0:
		prefix = "[DEBUG]"
	case 100:
		prefix = "[INFO]"
	default:
		prefix = "[INFO]"
	}

	if loglvl >= config.LogLevel {
		var str string
		for _, v := range msg {
			str = str + v
		}
		fmt.Printf("%-8s %s \n", prefix, str)
	}
}

func parseyaml() bool {
	var ok = true

	source, readerr := ioutil.ReadFile("config.yaml")

	if readerr != nil {
		ok = false
		consoleapp(4, readerr.Error())
		consoleapp(4, "Error while reading config file")
		return ok
	}

	parseerr := yaml.UnmarshalStrict(source, &config) //existing key override default
	if parseerr != nil {
		ok = false
		consoleapp(4, parseerr.Error())
		consoleapp(4, "Error while parsing config file")
		return ok
	}

	return ok
}

func printconfig() {
	if config.LogLevel == 0 {
		fmt.Println("Running-Config : ", config)
	}
}

func printcookie(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, c := range cookies {
		consoleapp(0, "cookie name : ", c.Name, " ; cookie value : ", c.Value)
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

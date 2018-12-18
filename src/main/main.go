package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/yaml.v2"
)

/* TODO:
READ-ONLY: empecher les requetes MOVE & DELETE
AMELIORER COOKIE (hash, ) raw champ ?
ENCRYPT PASSWORD
CONFIG INTERFACE
r.RemoteAddr ==> @ip:port => le port client peut changer vu qu'il est random
donc il est possible que on regénère un token pour rien(meme si le port est fixe pour une session et
que les client ne memorise les cookie souvent que pour la session en dav)
***
Le systeme de cookie reduit les risques mais n'est pas tres efficace,
on peut quand meme se connecter avec les id si on ne renvoie pas le cookie.
par contre il permet de reduire le traffix reseaux puisque il n'y a plus besoin
d'envoyer les ids a chaque requete (x2/3 du traffic)
*/

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
	logger(1, "Reading & Parsing Config")
	var ok bool
	ok = parseyaml()
	if !ok {
		logger(3, "Falling back to default config")
		config = defaultconf
	}

	printconfig() //config.LogLevel 0

	serverlistenon = config.ServerAddr + ":" + strconv.Itoa(config.ServerPort)

	davhandler = createdavhandler()

	logger(1, "Registrering Handlers") //on ajoute les handlers
	//http.Handle(config.DavPrefix, davhandler) //test without auth
	http.Handle(config.DavPrefix, authandler(davhandler))
	http.HandleFunc("/cookielist", printcookie)
	http.HandleFunc("/cookieset", setcookie)
}

func main() {
	initDB()
	testdb()
	logger(-1, "Starting Server")
	logger(-1, "Listening on", serverlistenon)
	fmt.Println("***********************************************************")
	http.ListenAndServe(serverlistenon, nil) //on demarre le serveur
}

func parseyaml() bool {
	var ok = true

	source, readerr := ioutil.ReadFile("config.yaml")

	if readerr != nil {
		ok = false
		logger(4, readerr.Error())
		logger(4, "Error while reading config file")
		return ok
	}

	parseerr := yaml.UnmarshalStrict(source, &config) //existing key override default
	if parseerr != nil {
		ok = false
		logger(4, parseerr.Error())
		logger(4, "Error while parsing config file")
		return ok
	}

	return ok
}

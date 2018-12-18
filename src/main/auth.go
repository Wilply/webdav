package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func authandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //handlerFunc AVEC UN R, s'occupe de passer les arguments a la func
		//1st test if cookie
		c, err := r.Cookie("Auth")
		if err == nil {
			ok, dbip, dbuser, dbtoken, dbexpire := getconnectionbyip(r.RemoteAddr)
			if ok && c.Value == dbtoken && r.RemoteAddr == dbip {
				if checkexpiretime(dbexpire) {
					next.ServeHTTP(w, r) //time before expiration > 60s
					return
				} else { //time before expiration < 60s
					http.SetCookie(w, newconnection(dbuser, r.RemoteAddr)) //overwrite cookie by generating new one and update db for another 20 min
					next.ServeHTTP(w, r)
					return
				}
			} else {
				delconnection(r.RemoteAddr) //delte invalide connection
				logger(1, "Invalid cookie for "+r.RemoteAddr+" try to re-authenticate")
			}
		} else {
			logger(1, "No cookie for "+r.RemoteAddr+" try to authenticate")
		}
		var user, pass, ok = r.BasicAuth()
		if !ok { // header authorization no set
			logger(2, "connection attempt from ", r.RemoteAddr, " set WWW-Authenticate header")
			w.Header().Set("WWW-Authenticate", `Basic realm="webdav acces, auth required" charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			if authenticate(user, pass) { //fonction d'authentification
				http.SetCookie(w, newconnection(user, r.RemoteAddr))
				logger(2, "User : ", user, " is connected from ", r.RemoteAddr)
			} else {
				logger(3, "Failed connexion attempt from ", r.RemoteAddr, " with user : ", user)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func authenticate(user, pass string) bool { //on modifiera ca plus tard
	ok, _, user_pw, _, _ := getuser(user) //if ok == true then user exist
	if ok && pass == user_pw {            //TODO: ENCRYPT
		return true
	}
	return false
}

func generatetoken() (token string) {
	btab := make([]byte, 10)       //genere un tableau de 8 byte de zero
	rand.Read(btab)                //remplace les zero par des nombre aleatoire
	return fmt.Sprintf("%x", btab) //print les nombre en base 16(%x) les un a cote des autre dans une seule chaine
}

func newconnection(username, ipaddr string) (c *http.Cookie) {
	token := generatetoken()
	expire := time.Now().Add(time.Second * 1200) // 20 min
	insertconnection(ipaddr, username, token, timetostring(expire))
	c = &http.Cookie{
		Name:    "Auth",
		Value:   token,
		Expires: expire,
	}
	return
}

func checkexpiretime(expiret string) bool {
	oldtime, _ := strconv.Atoi(expiret)
	dt := oldtime - timetoint(time.Now())
	if dt > 60 {
		return true
	}
	return false
}

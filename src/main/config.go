package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func gethashpassword(clearpass string) (hashedpass string) {
	bytepass := []byte(clearpass)
	pass, _ := bcrypt.GenerateFromPassword(bytepass, 12)
	hashedpass = string(pass)
	return
}

func matchpassword(hashedpass, clearpass string) bool {
	bytepass := []byte(clearpass)
	hspass := []byte(hashedpass)
	r := bcrypt.CompareHashAndPassword(hspass, bytepass)
	if r != nil {
		fmt.Println("wrong pass")
		return false
	}
	return true
}

func createuser(name, pass string) (ok bool) {
	okget, _, _, _, _ := getuser(name)
	if !okget {
		okinsert := insertuserDB(name, gethashpassword(pass), "", "")
		if okinsert {
			return true
		}
		logger(3, "Cannot create user "+name)
	}
	logger(3, "Cannot create user "+name+" user already exist")
	return false
}

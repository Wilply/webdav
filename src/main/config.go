package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashpassword(clearpass string) (hashedpass string) {
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

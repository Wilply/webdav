package main

import (
	"fmt"
	"log"
)

func logger(loglvl int, msg ...string) { //affiche les info si config.LogLevel > 1
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
		log.Printf("%-8s %s \n", prefix, str)
	}
}

func iferror(dept int, r error) {
	if r != nil {
		logger(dept, r.Error())
		if dept == 4 {
			log.Panic("The program encounter a problem and stopped")
		}
	}
}

func testr(r error) (ok bool) {
	if r != nil {
		ok = false
	} else {
		ok = true
	}
	return
}

func printconfig() {
	if config.LogLevel == 0 {
		fmt.Println("Running-Config : ", config)
	}
}

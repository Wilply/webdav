package main

import (
	"fmt"
	"strconv"
	"time"
)

func gettime() string {
	return fmt.Sprint(time.Now().Format("20060102150405")) //YYYYMMDDHHMMSS
}

func timeplushour() string {
	return fmt.Sprint(time.Now().Add(time.Second * 3600).Format("20060102150405")) //YYYYMMDDHHMMSS
}

func timeplusmin() string {
	return fmt.Sprint(time.Now().Add(time.Second * 60).Format("20060102150405")) //YYYYMMDDHHMMSS
}

func timeplus20min() string {
	return fmt.Sprint(time.Now().Add(time.Second * 20 * 60).Format("20060102150405")) //YYYYMMDDHHMMSS
}

func timetostring(t time.Time) string {
	return t.Format("20060102150405") //YYYYMMDDHHMMSS
}

func timetoint(t time.Time) (tin int) {
	tin, _ = strconv.Atoi(t.Format("20060102150405"))
	return
}

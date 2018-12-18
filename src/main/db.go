package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbrqlist = []string{
	"CREATE TABLE IF NOT EXISTS Users (user_login TEXT NOT NULL PRIMARY KEY,user_pass TEXT NOT NULL,user_rw TEXT,user_ro TEXT);",
	"CREATE TABLE IF NOT EXISTS Groups (group_name TEXT NOT NULL PRIMARY KEY,group_rw TEXT,group_ro TEXT);",
	"CREATE TABLE IF NOT EXISTS InGroup (user_login TEXT REFERENCES Users(user_login), group_name TEXT REFERENCES Groups(group_name));",
	"CREATE TABLE IF NOT EXISTS Connections (ipaddr TEXT PRIMARY KEY, user_login TEXT REFERENCES Users(user_login), token TEXT NOT NULL, expire TEXT NOT NULL);",
}

func testdb() {
	/*fmt.Println(insertuserDB("root", "azerty", "/", ""))
	listuser()
	insertgroupDB("admingrfkfk", "", "/dav")
	listgroup()
	insertingroupDB("root", "admingrp")
	listingroup()
	fmt.Println(getuser("root"))
	fmt.Println(getgroup("admingrfkfk"))
	fmt.Println(getgroupofuser("root"))
	fmt.Println(getusersingroup("admingrp"))
	fmt.Println(insertconnection("root", "CECIESTUNTOKEN", "ladate"))
	listconnection()
	fmt.Println(getconnection("root"))*/
	insertuserDB("root", "azerty", "/", "")
	insertgroupDB("admin", "", "")
	insertingroupDB("root", "admin")
}

func initDB() {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(4, r)
	for _, rqstr := range dbrqlist {
		rq, r := db.Prepare(rqstr)
		iferror(4, r)
		defer db.Close()
		rq.Exec()
	}
	db.Close()
}

func insertuserDB(username, password, rw, ro string) (ok bool) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	var rqstr string = "INSERT INTO Users VALUES (?,?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	_, r = rq.Exec(username, password, rw, ro)
	ok = testr(r)
	db.Close()
	return
}

func insertgroupDB(groupname, rw, ro string) (ok bool) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	var rqstr string = "INSERT INTO Groups VALUES (?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	_, r = rq.Exec(groupname, rw, ro)
	ok = testr(r)
	db.Close()
	return
}

func insertingroupDB(username, groupname string) (ok bool) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	//DELETE TO AVOID DUPLICATE, it's ugly but it's work and it's 4am
	var rqstrdel string = "DELETE FROM InGroup WHERE user_login = ? AND group_name = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(4, r)
	defer db.Close()
	rqdel.Exec(username, groupname)
	//INSERT
	var rqstr string = "INSERT INTO InGroup VALUES(?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	_, r = rq.Exec(username, groupname)
	ok = testr(r)
	db.Close()
	return
}

func listuser() {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Users")
	iferror(4, r)
	defer db.Close()
	var name, pass, rw, ro string
	fmt.Println("### LIST USER IN DB ###")
	for rows.Next() {
		rows.Scan(&name, &pass, &rw, &ro)
		fmt.Printf("# %-10s %s %-15s %s %-10s %s %-10s \n", name, " | ", pass, " | ", rw, " | ", ro)
	}
	fmt.Println("###")
}

func listgroup() {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Groups")
	iferror(4, r)
	defer db.Close()
	var name, rw, ro string
	fmt.Println("### LIST GROUP IN DB ###")
	for rows.Next() {
		rows.Scan(&name, &rw, &ro)
		fmt.Printf("# %-10s %s %-15s %s %-10s \n", name, " | ", rw, " | ", ro)
	}
	fmt.Println("###")
}

func listingroup() {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM InGroup")
	iferror(4, r)
	defer db.Close()
	var name, group string
	fmt.Println("### LIST INGROUP IN DB ###")
	for rows.Next() {
		rows.Scan(&name, &group)
		fmt.Printf("# %-10s %s %-10s \n", name, " | ", group)
	}
	fmt.Println("###")
}

func getuser(username string) (ok bool, login, pass, rw, ro string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Users WHERE user_login = ?", username)
	iferror(4, r)
	defer db.Close()
	r = row.Scan(&login, &pass, &rw, &ro)
	ok = testr(r)
	return
}

func getgroup(groupname string) (ok bool, name, rw, ro string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Groups WHERE group_name = ?", groupname)
	iferror(4, r)
	defer db.Close()
	r = row.Scan(&name, &rw, &ro)
	ok = testr(r)
	return
}

func getgroupofuser(username string) (ok bool, groups []string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT group_name FROM InGroup WHERE user_login = ?", username)
	iferror(4, r)
	defer db.Close()
	var tmpgrp string
	for rows.Next() {
		rows.Scan(&tmpgrp)
		groups = append(groups, tmpgrp)
	}
	ok = testr(r)
	return
}

func getusersingroup(groupname string) (ok bool, users []string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT user_login FROM InGroup WHERE group_name = ?", groupname)
	iferror(4, r)
	defer db.Close()
	var tmpusr string
	for rows.Next() {
		rows.Scan(&tmpusr)
		users = append(users, tmpusr)
	}
	ok = testr(r)
	return
}

func listconnection() {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Connections")
	iferror(4, r)
	defer db.Close()
	var ipaddr, login, token, expire string
	fmt.Println("### LIST CONNECTION IN DB ###")
	for rows.Next() {
		rows.Scan(&ipaddr, &login, &token, &expire)
		fmt.Printf("# %-12s %s %-10s %s %-20s %s %-15s \n", ipaddr, " | ", login, " | ", token, " | ", expire)
	}
	fmt.Println("###")
}

func getconnectionbyuser(login string) (ok bool, connection []string) { //DO NOT USE, can return more than one row
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Connections WHERE user_login = ?", login)
	iferror(4, r)
	defer db.Close()
	var tmpip, tmplogin, tmptoken, tmpexpire string
	for rows.Next() {
		rows.Scan(&tmpip, &tmplogin, &tmptoken, &tmpexpire)
		connection = append(connection, tmpip, tmplogin, tmptoken, tmpexpire)
	}
	ok = testr(r)
	return
}

func getconnectionbyip(ipaddr string) (ok bool, ip, user_login, token, expire string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Connections WHERE ipaddr = ?", ipaddr)
	iferror(4, r)
	defer db.Close()
	r = row.Scan(&ip, &user_login, &token, &expire)
	ok = testr(r)
	return
}

func insertconnection(ipaddr, login, token, expire string) (ok bool) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	//DELETE TO AVOID DUPLICATE, it's ugly but it's work and it's 4am
	var rqstrdel string = "DELETE FROM Connections WHERE ipaddr = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(4, r)
	defer db.Close()
	rqdel.Exec(ipaddr)
	//INSERT
	var rqstr string = "INSERT INTO Connections VALUES(?,?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	_, r = rq.Exec(ipaddr, login, token, expire)
	ok = testr(r)
	db.Close()
	return
}

func delconnection(ipaddr string) (ok bool) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	var rqstrdel string = "DELETE FROM Connections WHERE ipaddr = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(4, r)
	defer db.Close()
	rqdel.Exec(ipaddr)
	ok = testr(r)
	db.Close()
	return
}

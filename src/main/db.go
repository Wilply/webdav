package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

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
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(4, r)
	for _, rqstr := range dbrqlist {
		rq, r := db.Prepare(rqstr)
		iferror(4, r)
		defer db.Close()
		rq.Exec()
	}
	db.Close()
}

// #################
// ##### USERS #####
// #################

func insertuserDB(username, password, rw, ro string) (ok bool) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstr string = "INSERT INTO Users VALUES (?,?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(username, password, rw, ro)
	ok = testr(r)
	if ok {
		logger(1, "Insert user ", username, "succesful")
	}
	db.Close()
	return
}

func getuser(username string) (ok bool, login, pass, rw, ro string) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Users WHERE user_login = ?", username)
	iferror(3, r)
	defer db.Close()
	r = row.Scan(&login, &pass, &rw, &ro)
	ok = testr(r)
	return
}

func deleteuser(username string) (ok bool) {
	userexist, _, _, _, _ := getuser(username)
	if !userexist {
		logger(3, "cannot delete user : "+username+" , user does not exist")
		ok = false
		return
	}
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstr string = "DELETE FROM Users WHERE user_login = ?;"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(username)
	ok = testr(r)
	if ok {
		logger(1, "Delete user ", username, "succesful")
	}
	db.Close()
	return
}

func updateuserpassword(username, password string) (ok bool) {
	userexist, _, _, _, _ := getuser(username)
	if !userexist {
		logger(3, "cannot update user : "+username+" , user does not exist")
		ok = false
		return
	}
	if password != "" {
		db, r := sql.Open("sqlite3", config.DBfile)
		iferror(3, r)
		var rqstr string = "UPDATE Users SET user_pass = ? WHERE user_login = ?;"
		rq, r := db.Prepare(rqstr)
		iferror(3, r)
		defer db.Close()
		_, r = rq.Exec(password, username)
		ok = testr(r)
		if ok {
			logger(1, "Update password for  user ", username, "succesful")
		}
		db.Close()
	} else {
		logger(3, "cannot update user, new password too short")
	}
	return
}

func updateuserro(username, ro string) (ok bool) {
	userexist, _, _, _, _ := getuser(username)
	if !userexist {
		logger(3, "cannot update user : "+username+" , user does not exist")
		ok = false
		return
	}
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstr string = "UPDATE Users SET user_ro = ? WHERE user_login = ?;"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(ro, username)
	ok = testr(r)
	if ok {
		logger(1, "Update ro for user ", username, "succesful")
	}
	db.Close()
	return
}

func updateuserrw(username, rw string) (ok bool) {
	userexist, _, _, _, _ := getuser(username)
	if !userexist {
		logger(3, "cannot update user : "+username+" , user does not exist")
		ok = false
		return
	}
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstr string = "UPDATE Users SET user_rw = ? WHERE user_login = ?;"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(rw, username)
	ok = testr(r)
	if ok {
		logger(1, "Update rw for user ", username, "succesful")
	}
	db.Close()
	return
}

// ##################
// ##### GROUPS #####
// ##################

func insertgroupDB(groupname, rw, ro string) (ok bool) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstr string = "INSERT INTO Groups VALUES (?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(groupname, rw, ro)
	ok = testr(r)
	if ok {
		logger(1, "Insert ", groupname, "succesful")
	}
	db.Close()
	return
}

func getgroup(groupname string) (ok bool, name, rw, ro string) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Groups WHERE group_name = ?", groupname)
	iferror(3, r)
	defer db.Close()
	r = row.Scan(&name, &rw, &ro)
	ok = testr(r)
	return
}

// ####################
// ##### INGROUPS #####
// ####################

func insertingroupDB(username, groupname string) (ok bool) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	//DELETE TO AVOID DUPLICATE, it's ugly but it's work and it's 4am
	var rqstrdel string = "DELETE FROM InGroup WHERE user_login = ? AND group_name = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(3, r)
	defer db.Close()
	rqdel.Exec(username, groupname)
	//INSERT
	var rqstr string = "INSERT INTO InGroup VALUES(?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(username, groupname)
	ok = testr(r)
	if ok {
		logger(1, "Insert InGroup ", username, " : ", groupname, "succesful")
	}
	db.Close()
	return
}

func getusersingroup(groupname string) (ok bool, users []string) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT user_login FROM InGroup WHERE group_name = ?", groupname)
	iferror(3, r)
	defer db.Close()
	var tmpusr string
	for rows.Next() {
		rows.Scan(&tmpusr)
		users = append(users, tmpusr)
	}
	ok = testr(r)
	return
}

func getgroupofuser(username string) (ok bool, groups []string) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT group_name FROM InGroup WHERE user_login = ?", username)
	iferror(3, r)
	defer db.Close()
	var tmpgrp string
	for rows.Next() {
		rows.Scan(&tmpgrp)
		groups = append(groups, tmpgrp)
	}
	ok = testr(r)
	return
}

// #######################
// ##### CONNECTIONS #####
// #######################

func insertconnection(ipaddr, login, token, expire string) (ok bool) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	//DELETE TO AVOID DUPLICATE, it's ugly but it's work and it's 4am
	var rqstrdel string = "DELETE FROM Connections WHERE ipaddr = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(3, r)
	defer db.Close()
	rqdel.Exec(ipaddr)
	//INSERT
	var rqstr string = "INSERT INTO Connections VALUES(?,?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(3, r)
	defer db.Close()
	_, r = rq.Exec(ipaddr, login, token, expire)
	ok = testr(r)
	db.Close()
	return
}

func getconnectionbyuser(login string) (ok bool, connection []string) { //DO NOT USE, can return more than one row
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Connections WHERE user_login = ?", login)
	iferror(3, r)
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
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	row := db.QueryRow("SELECT * FROM Connections WHERE ipaddr = ?", ipaddr)
	iferror(3, r)
	defer db.Close()
	r = row.Scan(&ip, &user_login, &token, &expire)
	ok = testr(r)
	return
}

func delconnection(ipaddr string) (ok bool) {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	var rqstrdel string = "DELETE FROM Connections WHERE ipaddr = ?;"
	rqdel, r := db.Prepare(rqstrdel)
	iferror(3, r)
	defer db.Close()
	rqdel.Exec(ipaddr)
	ok = testr(r)
	db.Close()
	return
}

func clearconnection() {
	for true {
		time.Sleep(2 * time.Minute)
		nowtime := timetoint(time.Now())
		db, r := sql.Open("sqlite3", config.DBfile)
		iferror(3, r)
		rows, r := db.Query("SELECT ipaddr, expire FROM Connections")
		iferror(3, r)
		defer db.Close()
		var ipaddr, expire string
		for rows.Next() {
			rows.Scan(&ipaddr, &expire)
			oldtime, _ := strconv.Atoi(expire)
			if (oldtime - nowtime) < 0 {
				delconnection(ipaddr)
			}
		}
	}
}

// ##### TESTS ####

func listuser() {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Users")
	iferror(3, r)
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
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Groups")
	iferror(3, r)
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
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM InGroup")
	iferror(3, r)
	defer db.Close()
	var name, group string
	fmt.Println("### LIST INGROUP IN DB ###")
	for rows.Next() {
		rows.Scan(&name, &group)
		fmt.Printf("# %-10s %s %-10s \n", name, " | ", group)
	}
	fmt.Println("###")
}

func listconnection() {
	db, r := sql.Open("sqlite3", config.DBfile)
	iferror(3, r)
	rows, r := db.Query("SELECT * FROM Connections")
	iferror(3, r)
	defer db.Close()
	var ipaddr, login, token, expire string
	fmt.Println("### LIST CONNECTION IN DB ###")
	for rows.Next() {
		rows.Scan(&ipaddr, &login, &token, &expire)
		fmt.Printf("# %-12s %s %-10s %s %-20s %s %-15s \n", ipaddr, " | ", login, " | ", token, " | ", expire)
	}
	fmt.Println("###")
}

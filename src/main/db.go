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
}

func testdb() {
	insertuserDB("root", "azerty", "/", "")
	listuser()
	insertgroupDB("admingrfkfk", "", "/dav")
	listgroup()
	insertingroupDB("root", "admingrp")
	listingroup()
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

func insertuserDB(username, password, rw, ro string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	var rqstr string = "INSERT INTO Users VALUES (?,?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	rq.Exec(username, password, rw, ro)
	db.Close()
}

func insertgroupDB(groupname, rw, ro string) {
	db, r := sql.Open("sqlite3", "./foo.db")
	iferror(3, r)
	var rqstr string = "INSERT INTO Groups VALUES (?,?,?);"
	rq, r := db.Prepare(rqstr)
	iferror(4, r)
	defer db.Close()
	rq.Exec(groupname, rw, ro)
	db.Close()
}

func insertingroupDB(username, groupname string) {
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
	rq.Exec(username, groupname)
	db.Close()
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

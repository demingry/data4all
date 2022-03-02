package main

import "database/sql"

func HeaderReserve() {

	dsn := `username:password@tcp(127.0.0.1:3306)/dbname`
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

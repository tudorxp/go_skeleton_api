package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"errors"
	"time"
)


const (
	CONN_TIMEOUT = 5 * time.Second 	// Timeout for connecting to the database
	CONN_LIFETIME = 1 * time.Hour  // Max connection lifetime 
)

// db_init() sets up a connection to the database on the connstring/URL supplied
// it returns a *sql.DB object or an error
// caller is responsible for closing the *sql.DB
func db_init(connstr string, max_open_conns int, max_idle_conns int) (db *sql.DB, err error) {

	timeout_channel := time.After(CONN_TIMEOUT)
	c:=make(chan struct{})

	go func() {
		defer func(){
			c <- struct{}{}
		}()
		db, err = sql.Open("postgres",connstr)
		if err!=nil {
			return
		}
		db.SetConnMaxLifetime(CONN_LIFETIME)
		db.SetMaxIdleConns(max_idle_conns)
		db.SetMaxOpenConns(max_open_conns)
		err = db.Ping()
		if err != nil {
			return 
		}		
	}()

	select {
	case <- c:
		return db, err
	case <- timeout_channel:
		return nil, errors.New("Timeout reached while trying to connect to database")
	}

}
package main

import (
	"log"
	"database/sql"
	"net/http"
	"io"
	"time"
	"context"
)

type Server struct {
	log *log.Logger
	db *sql.DB
	mux *http.ServeMux
	httpserver *http.Server
}

const (
	LOGGER_FORMAT = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	LOGGER_PREFIX = ""
	HTTP_READHEADERTIMEOUT = 20 *time.Second
	SHUTDOWN_TIMEOUT = 20*time.Second
)

// Configures and returns a server object
func NewServer(db_connstr string, db_maxopenconns int, db_maxidleconns int, logger_out io.Writer, bind_address string) (*Server, error) {

	// DB
	db, err := db_init(db_connstr, db_maxopenconns, db_maxidleconns)
	if err!= nil {
		return nil, err
	}
	
	// Log
	log:=log.New(logger_out,LOGGER_PREFIX,LOGGER_FORMAT)

	// Mux
	mux:= http.NewServeMux()

	// HTTP Server
	httpserver := &http.Server{
		Addr: bind_address,
		Handler: mux,
		ErrorLog: log,
		ReadHeaderTimeout: HTTP_READHEADERTIMEOUT,
	}

	return &Server{
		db: db,
		log: log,
		mux: mux,
		httpserver: httpserver,
	}, nil

}

// Initialises and starts a server
func (s *Server) Start() {
	s.log.Print("Starting server.")

	s.addRoutes()

	// trap exit
	go func() {
		err := s.httpserver.ListenAndServe()
		s.log.Printf("HTTP server exited with: %s", err)

	}()

}


// Shuts down the HTTP server
func (s *Server) Stop() {

	c, _:=context.WithTimeout(context.Background(),SHUTDOWN_TIMEOUT)
	err:=s.httpserver.Shutdown(c)
	if err!=nil {
		s.log.Printf("HTTP server shutdown error: %s",err)
	} else {
		s.log.Printf("HTTP server shutdown successful.")
	}

	err = s.db.Close()
	if err!= nil {
		s.log.Printf("DB close error: %s", err)
	} else {
		s.log.Printf("Closed DB connections.")
	}

}
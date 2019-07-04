package main

import (
	"fmt"
	"net/http"
	"path"
	"encoding/json"
)

// Configures the mux for HTTP routing
func (s *Server) addRoutes() {
	s.mux.HandleFunc("/hello",s.helloHandler())
	s.mux.HandleFunc("/sayHi",s.sayHiHandler1)
	s.mux.HandleFunc("/sayHi/",s.sayHiHandler2)
	s.mux.HandleFunc("/listUsers",s.listUsersHandler())
}

// Defining and returning a HandlerFunc enables chaining by passing a "next" parameter
func (s *Server) helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,"Hello, hello.")
	}
}

func (s *Server) sayHiHandler1(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Error: Method not allowed.")
		return
	}
	name:=r.FormValue("name")
	if name=="" {
		http.Error(w,"Error: You must provide a name.",http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w,"Hi, %s.",name)
}

func (s *Server) sayHiHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Error: Method not allowed.")
		return
	}
	
	dir, rest:=path.Split(r.URL.Path)
	if dir != "/sayHi/" { // Possibly multiple slashes
		http.Error(w,"Error: Bad URL",http.StatusBadRequest)
		return		
	}
	if rest == "" {
		http.Error(w,"Error: You must provide a name.",http.StatusBadRequest)
		return		
	}
	fmt.Fprintf(w,"Hi, %s.",rest)
}

// DB access
func (s *Server) listUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err:=ListUsers(s.db)
		if err!=nil {
			http.Error(w,fmt.Sprintf("Error getting list of users: %s",err),http.StatusInternalServerError)
			return
		}
	
		b, err:=json.Marshal(users)
		if err!=nil {
			http.Error(w,fmt.Sprintf("Error marshalling into JSON: %s",err),http.StatusInternalServerError)
			return
		}		
		w.Write(b)

	}
}

package main

import (
	"testing"
	"io/ioutil"
)

func testError(t *testing.T, err error) {
	if err!=nil {
		t.Error(err)
	}
}

func testFatal(t *testing.T, err error) {
	if err!=nil {
		t.Fatal(err)
	}
}


func TestNewServer (t *testing.T) {
	s, err := NewServer("postgresql://top:top@localhost:5432/top?sslmode=disable", 10, 2, ioutil.Discard , ":9123")
	testFatal(t,err)
	switch {
	case s.db == nil:
		t.Error("couldn't create DB")
	case s.log == nil:
		t.Error("couldn't create logger")
	case s.mux == nil:
		t.Error("couldn't create mux")
	case s.httpserver == nil:
		t.Error("couldn't create HTTP server")
	}
		
}


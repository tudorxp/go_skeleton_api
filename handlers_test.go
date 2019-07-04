package main

import (
	"os"
	"testing"
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"net/http"
	"strings"
)

var ts *Server

const (
	TESTDB="postgresql://top:top@localhost:5432/top?sslmode=disable"
)

// Init server & DB before tests 
func TestMain(m *testing.M) {
	flag.Parse()
	var err error
	s, err := NewServer(TESTDB, 10, 2, ioutil.Discard , ":9123")
	if err!=nil {
		panic(err)
	}
	ts=s
	s.Start()
	os.Exit(m.Run())
}

func TestHelloHandler(t *testing.T) {

	req:=httptest.NewRequest("GET","http://test.com/hello",nil)
	expected:="Hello, hello."
	w:=httptest.NewRecorder()
	ts.helloHandler()(w, req)
	resp:=w.Result()
	got, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("returned HTTP status: %d",resp.StatusCode)
	}
	if string(got) != "Hello, hello." {
		t.Errorf("expected: %s, got: %s",expected, got)
	}
}

func TestSayHiHandler1(t *testing.T) {

	tests := []struct { 
		method string
		input string
		statuscode int
		output string
	}{
		{ "GET", "", 405, "Error: Method not allowed.\n" },
		{ "POST", "", 400, "Error: You must provide a name.\n" },
		{ "POST", "Foo", 200, "Hi, Foo." },
		{ "POST", "Bar", 200, "Hi, Bar." },
	}

	for _, test := range tests {
		var body string
		if test.input=="" {
			body=""
		} else {
			body="name="+test.input
		}
		req:=httptest.NewRequest(test.method,"http://test.com/sayHi",strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w:=httptest.NewRecorder()
		ts.sayHiHandler1(w, req)
		resp:=w.Result()
		got, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != test.statuscode {
			t.Errorf("expected HTTP status code: %d, got: %d",test.statuscode, resp.StatusCode)
		}
		if string(got) != test.output {
			t.Errorf("expected: %s, got: %s",test.output, got)
		}

	}
}
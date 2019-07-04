package main

import (
	// "fmt"
	"flag"
	"os"
	"os/signal"
)

func main() {

	var LOG_DESTINATION = os.Stdin // or for web scale, ioutil.Discard	

	var (
		db_connstr = flag.String("db_connstr","postgresql://top:top@localhost:5432/top?sslmode=disable","URL or connstr to use for DB connection")
		db_maxopenconns = flag.Int("db_maxopenconns",10,"Size of the DB connection pool")
		db_maxidleconns = flag.Int("db_maxidleconns",2,"Max idle connections in DB connection pool")	
		bind_address = flag.String("bind_address","localhost:8088","HTTP server bind address")		
	)

	flag.Parse()

	s, err := NewServer(*db_connstr, *db_maxopenconns, *db_maxidleconns, LOG_DESTINATION, *bind_address)
	if err!=nil {
		panic(err)
	}

	s.Start()

	// trap SIGINT to initiate server shutdown
	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)

	go func() { // shutdown
		<- signal_channel
		s.Stop()
		close(signal_channel)
		os.Exit(0)
	}()

	select {}
}

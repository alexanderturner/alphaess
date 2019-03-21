package main

import (
	"log"
	"os"
)

// DBURL, DBUSER, DBPASS, DBNAME, ESSUSER, ESSPASS, ESSSN

func main() {
	//test ENV Vars
	if os.Getenv("DBURL") == "" {
		log.Fatal("DBURL not set")
	}
	if os.Getenv("DBUSER") == "" {
		log.Fatal("DBUSER not set")
	}
	if os.Getenv("DBPASS") == "" {
		log.Fatal("DBPASS not set")
	}
	if os.Getenv("DBNAME") == "" {
		log.Fatal("DBNAME not set")
	}
	if os.Getenv("ESSUSER") == "" {
		log.Fatal("ESSUSER not set")
	}
	if os.Getenv("ESSPASS") == "" {
		log.Fatal("ESSPASS not set")
	}
	if os.Getenv("ESSSN") == "" {
		log.Fatal("ESSSN not set")
	}

	done := make(chan bool)
	a := new(App)
	a.DB = influxDBClient()
	a.apiPoller()

	<-done

}

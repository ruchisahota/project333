package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func cassInit() {
	var err error

	log.Println("Cassandra Details ", os.Getenv("CASS_IP"), os.Getenv("CASS_PORT"), os.Getenv("CASS_KEYSPACE"))

	cluster := gocql.NewCluster(os.Getenv("CASS_IP") + ":" + os.Getenv("CASS_PORT"))
	cluster.Keyspace = os.Getenv("CASS_KEYSPACE")
	//cluster.ProtoVersion = 3
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra connection established...")
}

func closeCass() {
	time.Sleep(2 * time.Second)
	Session.Close()
	fmt.Println("cassandra connection closed...")
}

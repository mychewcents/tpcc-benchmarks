package common

import (
	"encoding/xml"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/config"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type CassandraSession struct {
	ReadSession  *gocql.Session
	WriteSession *gocql.Session
}

func MakeCassandraSession(path string) *CassandraSession {
	cassandraConfig := makeCassandraConfig(path)

	readCluster := gocql.NewCluster(cassandraConfig.Hosts...)
	readCluster.Keyspace = "cassandra"
	readCluster.Timeout = time.Minute * 2
	readCluster.NumConns = 10
	if strings.ToUpper(cassandraConfig.ReadConsistency) == "ONE" {
		readCluster.Consistency = gocql.LocalOne
	} else {
		readCluster.Consistency = gocql.LocalQuorum
	}
	readCluster.SerialConsistency = gocql.LocalSerial
	readSession, err := readCluster.CreateSession()
	if err != nil {
		panic("error creating cassandra session for read")
	}

	writeCluster := gocql.NewCluster(cassandraConfig.Hosts...)
	writeCluster.Keyspace = "cassandra"
	writeCluster.Timeout = time.Minute * 2
	readCluster.NumConns = 10
	if strings.ToUpper(cassandraConfig.WriteConsistency) == "ALL" {
		writeCluster.Consistency = gocql.All
	} else {
		writeCluster.Consistency = gocql.LocalQuorum
	}
	writeCluster.SerialConsistency = gocql.LocalSerial
	writeSession, err := writeCluster.CreateSession()
	if err != nil {
		panic("error creating cassandra session for write")
	}

	return &CassandraSession{
		ReadSession:  readSession,
		WriteSession: writeSession,
	}
}

func makeCassandraConfig(path string) *config.CassandraConfig {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var cassandraConfig config.CassandraConfig
	err = xml.Unmarshal(byteValue, &cassandraConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &cassandraConfig
}

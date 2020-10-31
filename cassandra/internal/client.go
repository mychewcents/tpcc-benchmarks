package cassandra_client

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/config"
	"github.com/mychewcents/ddbms-project/cassandra/internal/router"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var experimentId, clientId int

func init() {
	if len(os.Args) < 3 {
		panic("need to supply experimentId and clientId")
	}

	experimentId, _ = strconv.Atoi(os.Args[1])
	clientId, _ = strconv.Atoi(os.Args[2])

	fileName := fmt.Sprintf("log/logs_exp_%v_client%v", experimentId, clientId)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func Start() {
	log.Printf("Starting Client %v", clientId)
	defer log.Printf("Stopping Client %v", clientId)
	time.Sleep(10 * time.Second)

	cassandraSession := makeCassandraSession()
	reader := bufio.NewReader(os.Stdin)

	r := router.NewTransactionRouter(cassandraSession, reader)

	text, _ := reader.ReadString('\n')

	for text != "" {
		r.HandleCommand(text)
		text, _ = reader.ReadString('\n')
	}
}

func makeCassandraSession() *common.CassandraSession {
	cassandraConfig := makeCassandraConfig()

	readCluster := gocql.NewCluster(cassandraConfig.Hosts...)
	readCluster.Keyspace = cassandraConfig.Keyspace
	readCluster.Timeout = time.Minute * 2
	readCluster.NumConns = 10
	if strings.ToUpper(cassandraConfig.ReadConsistency) == "ONE" {
		readCluster.Consistency = gocql.One
	} else {
		readCluster.Consistency = gocql.Quorum
	}
	readSession, _ := readCluster.CreateSession()

	writeCluster := gocql.NewCluster(cassandraConfig.Hosts...)
	writeCluster.Keyspace = cassandraConfig.Keyspace
	writeCluster.Timeout = time.Minute * 2
	readCluster.NumConns = 10
	if strings.ToUpper(cassandraConfig.WriteConsistency) == "ONE" {
		writeCluster.Consistency = gocql.One
	} else {
		writeCluster.Consistency = gocql.Quorum
	}
	writeSession, _ := writeCluster.CreateSession()

	return &common.CassandraSession{
		ReadSession:  readSession,
		WriteSession: writeSession,
	}
}

func makeCassandraConfig() *config.CassandraConfig {
	xmlFile, err := os.Open("configs/local/cassandra-config.xml")
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

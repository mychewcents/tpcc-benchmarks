package cassandra_client

import (
	"bufio"
	"fmt"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/router"
	"log"
	"os"
	"strconv"
	"time"
)

var experimentId, clientId int

func Start() {
	if len(os.Args) < 4 {
		panic("need to supply experimentId, clientId and path to config")
	}

	experimentId, _ = strconv.Atoi(os.Args[1])
	clientId, _ = strconv.Atoi(os.Args[2])

	fileName := fmt.Sprintf("log/logs_exp_%v_client_%v", experimentId, clientId)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Printf("Starting Client %v", clientId)
	defer log.Printf("Stopping Client %v", clientId)

	time.Sleep(10 * time.Second)

	cassandraSession := common.MakeCassandraSession(os.Args[3])
	reader := bufio.NewReader(os.Stdin)

	r := router.NewTransactionRouter(cassandraSession, reader)
	m := controller.NewPerformanceMonitorController()

	for {
		start := time.Now()

		text, _ := reader.ReadString('\n')
		if text == "" {
			break
		}
		r.HandleCommand(text)

		end := time.Now()
		m.StoreLatency(int(end.Sub(start).Milliseconds()))
	}
	m.StorePerformanceMetrics("results/metrics", experimentId, clientId)
}

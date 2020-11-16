package cassandra_client

import (
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/controller"
	"os"
	"strconv"
)

func StoreDatabaseState() {
	if len(os.Args) < 3 {
		panic("need to supply experimentNo and path to config")
	}
	experimentId, _ = strconv.Atoi(os.Args[1])

	cassandraSession := common.MakeCassandraSession(os.Args[2])

	c := controller.NewDatabaseStateController(cassandraSession)
	c.SaveDatabaseState("results/dbstate", experimentId)
}

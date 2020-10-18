package cassandra_client

import (
	"bufio"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/router"
	"os"
)

func Start() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "cassandra"
	cluster.Consistency = gocql.One

	reader := bufio.NewReader(os.Stdin)
	r := router.NewTransactionRouter(cluster, reader)

	text, _ := reader.ReadString('\n')

	for text != "" {
		r.HandleCommand(text)
		text, _ = reader.ReadString('\n')
	}
}

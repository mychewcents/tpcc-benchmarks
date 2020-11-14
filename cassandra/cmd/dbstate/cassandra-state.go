package main

import cassandra_client "github.com/mychewcents/tpcc-benchmarks/cassandra/internal"

func main() {
	cassandra_client.StoreDatabaseState()
}

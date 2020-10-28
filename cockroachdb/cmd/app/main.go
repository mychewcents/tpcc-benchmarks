package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"
	"strings"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

var (
	experiment  = flag.Int("exp", 0, "Experiment Number")
	client      = flag.Int("client", 0, "Client Number")
	connPtr     = flag.String("host", "localhost", "URL / IP of the DB Server")
	portPtr     = flag.Int("port", 26257, "Port to contact the server's CDB Service")
	dbPtr       = flag.String("database", "defaultdb", "Database to connect")
	usernamePtr = flag.String("username", "root", "Username to connect with")
)

func init() {
	var err error
	flag.Parse()

	if *experiment == 0 || *client == 0 {
		panic("Provide Experiment and Client number to proceed")
	}
	db, err = cdbconn.CreateConnection(*connPtr, *portPtr, *dbPtr, *usernamePtr)
	if err != nil {
		panic(err)
	}
}

func main() {
	// fmt.Println(*experiment, *client)
	var txArgs []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		caller.ProcessRequest(db, scanner, txArgs)
	}
}

package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"
	"strings"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
)

var db *sql.DB

var (
	connPtr     = flag.String("host", "localhost", "URL / IP of the DB Server")
	portPtr     = flag.Int("port", 25267, "Port to contact the server's CDB Service")
	dbPtr       = flag.String("database", "defaultdb", "Database to connect")
	usernamePtr = flag.String("username", "root", "Username to connect with")
)

// func init() {
// 	var err error
// 	flag.Parse()

// 	db, err = cdbconn.CreateConnection(*connPtr, *portPtr, *dbPtr, *usernamePtr)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	var txArgs []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		caller.ProcessRequest(db, scanner, txArgs)
	}
	// for true {

	// 	var transaction_type byte
	// 	_, err := fmt.Scanf("%c", &transaction_type)

	// 	if err != nil {
	// 		if err.Error() == "EOF" {
	// 			fmt.Println("Read EOF")
	// 		} else {
	// 			fmt.Println(err)
	// 		}
	// 		break
	// 	}

	// 	switch transaction_type {
	// 	case 'N':
	// 	case 'P':
	// 		var warehouseId, districtId, customerId int
	// 		var amount float64
	// 		fmt.Scanf(",%d,%d,%d,%f", &warehouseId, &districtId, &customerId, &amount)
	// 		payment.ProcessTransaction(db, warehouseId, districtId, customerId, amount)
	// 		break
	// 	case 'D':
	// 	case 'O':
	// 		var warehouseId, districtId, customerId int
	// 		fmt.Scanf(",%d,%d,%d", &warehouseId, &districtId, &customerId)
	// 		order_status.ProcessTransaction(db, warehouseId, districtId, customerId)
	// 		break
	// 	case 'S':
	// 	case 'I':
	// 	case 'T':
	// 		top_balance.ProcessTransaction(db)
	// 		break
	// 	case 'R':

	// 	}
	// }
}

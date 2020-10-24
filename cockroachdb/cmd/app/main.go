package main

import (
	"database/sql"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

func init() {
	var err error
	db, err = cdbconn.CreateConnection("localhost", "26257", "defaultdb", "root")
	if err != nil {
		panic(err)
	}
}

func main() {
	// itemIDs := []int{
	// 	68195, 26567, 4114, 69343, 1836,
	// 	89294, 31975, 80487, 32423, 93383,
	// 	46279, 98511, 69351, 4679, 96983,
	// }
	// supplierIDs := []int{
	// 	1, 1, 1, 1, 1,
	// 	1, 1, 1, 1, 1,
	// 	1, 1, 1, 1, 1,
	// }
	// quantitys := []int{
	// 	1, 5, 7, 3, 1,
	// 	3, 6, 10, 5, 9,
	// 	9, 8, 4, 6, 2,
	// }

	// neworder.ProcessTransaction(db, 1, 1, 1279, 15, itemIDs, supplierIDs, quantitys)

	// popularitem.ProcessTransaction(db, 1, 1, 1000)
	// stocklevel.ProcessTransaction(db, 1, 1, 15, 1000)
	// for true {

	// var transaction_type byte
	// _, err := fmt.Scanf("%c", &transaction_type)

	// if err != nil {
	// 	if err.Error() == "EOF" {
	// 		fmt.Println("Read EOF")
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// 	break
	// }

	// switch (transaction_type) {
	// case 'N':
	// 	var c_id, d_id, w_id, m int
	// 	fmt.Scanf("%d %d %d %d", &c_id, &d_id, &w_id, &m)
	// 	var ol_i_id = make([]int64, m)
	// 	var ol_supply_w_id = make([]int64, m)
	// 	var ol_quantity = make([]float32, m)
	// 	for i := 0; i < m; i++ {
	// 		fmt.Scanf("%d %d %f", &ol_i_id[i], &ol_supply_w_id[i], &ol_quantity[i])
	// 	}
	// 	var new_order_input model.NewOrderInput
	// 	(&new_order_input).InitNewOrderInput(c_id, d_id, w_id, m, ol_i_id, ol_supply_w_id, ol_quantity)
	// 	executors.NewOrderExecutor(new_order_input)
	// 	new_order_input.PrintNewOrderInput()
	// 	break
	// 	case 'P':
	// 		var c_w_id, c_d_id, c_id int64
	// 		var payment float32
	// 		fmt.Scanf("%d %d %d %f", &c_w_id, &c_d_id, &c_id, &payment)

	// 		var payment_input model.PaymentInput
	// 		(&payment_input).InitPaymentInput(c_w_id, c_d_id, c_id, payment)
	// 		executors.PaymentExecutor(payment_input)
	// 		payment_input.PrintPaymentInput()
	// 		break
	// 	case 'D':
	// 	case 'O':
	// 	case 'S':
	// 	case 'I':
	// 	case 'T':
	// 	case 'R':

	// 	}
	// }
}

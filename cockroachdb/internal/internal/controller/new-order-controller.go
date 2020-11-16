package controller

import (
	"bufio"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// NewOrderControllerImpl provides the interface to call the service
type newOrderControllerImpl struct {
	s services.NewOrderService
}

// CreateNewOrderController get the new controller to execute the New Order Transaction
func CreateNewOrderController(db *sql.DB) handler.NewTransactionController {
	return &newOrderControllerImpl{
		s: services.CreateNewOrderService(db),
	}
}

// HandleTransaction performs the transaction and returns the execution result in Boolean
func (noc *newOrderControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {
	cID, _ := strconv.Atoi(args[0])
	wID, _ := strconv.Atoi(args[1])
	dID, _ := strconv.Atoi(args[2])
	numOfOrderLineItems, _ := strconv.Atoi(args[3])

	newOrderLines, isOrderLocal, totalUniqueItems := readAndPrepareOrderLineItems(scanner, numOfOrderLineItems, wID)

	n := &models.NewOrder{
		CustomerID:        cID,
		WarehouseID:       wID,
		DistrictID:        dID,
		IsOrderLocal:      isOrderLocal,
		UniqueItems:       totalUniqueItems,
		NewOrderLineItems: newOrderLines,
	}

	_, err := noc.s.ProcessTransaction(n)
	if err != nil {
		log.Printf("error found in the new order transaction. Err: %v", err)
		return false
	}

	return true
}

func readAndPrepareOrderLineItems(scanner *bufio.Scanner, numOfItems, warehouseID int) (orderLineItems map[int]*models.NewOrderOrderLineItem, isOrderLocal, totalUniqueOrderItems int) {
	orderLineItems = make(map[int]*models.NewOrderOrderLineItem)
	isOrderLocal = 1

	var id, supplier, quantity, remote int

	for i := 0; i < numOfItems; i++ {
		if scanner.Scan() {
			args := strings.Split(scanner.Text(), ",")

			id, _ = strconv.Atoi(args[0])
			supplier, _ = strconv.Atoi(args[1])
			quantity, _ = strconv.Atoi(args[2])

			if supplier != warehouseID {
				remote = 1
				if isOrderLocal == 1 {
					isOrderLocal = 0
				}
			} else {
				remote = 0
			}

			if _, ok := orderLineItems[id]; ok {
				orderLineItems[id].Quantity += quantity
			} else {
				orderLineItems[id] = &models.NewOrderOrderLineItem{
					SupplierWarehouseID: supplier,
					Quantity:            quantity,
					IsRemote:            remote,
				}
				totalUniqueOrderItems++
			}
		}
	}

	return
}

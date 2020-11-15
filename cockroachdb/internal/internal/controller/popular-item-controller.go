package controller

import (
	"bufio"
	"database/sql"
	"log"
	"strconv"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

type popularItemControllerImpl struct {
	s services.PopularItemService
}

// CreatePopularItemController creates a new controller to handle Popular Item transactions
func CreatePopularItemController(db *sql.DB) handler.NewTransactionController {
	return &popularItemControllerImpl{
		s: services.CreateNewPopularItemService(db),
	}
}

func (pic *popularItemControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {
	wID, _ := strconv.Atoi(args[0])
	dID, _ := strconv.Atoi(args[1])
	lastNOrders, _ := strconv.Atoi(args[2])

	pi := &models.PopularItem{
		WarehouseID: wID,
		DistrictID:  dID,
		LastNOrders: lastNOrders,
	}

	_, err := pic.s.ProcessTransaction(pi)
	if err != nil {
		log.Println("error occurred in executing the popular item transaction. Err: %v", err)
		return false
	}

	return true
}

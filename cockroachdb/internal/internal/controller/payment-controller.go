package controller

import (
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

type paymentControllerImpl struct {
	s services.PaymentService
}

func CreateNewPaymentController(db *sql.DB) handler.NewTransactionController {
	return &paymentControllerImpl{
		s: services.CreateNewPaymentService(db)
	}
}

func (pc *paymentControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {
	wID, _ := strconv.Atoi(args[0])
	dID, _ := strconv.Atoi(args[1])
	cID, _ := strconv.Atoi(args[2])
	paymentAmt, _ := strconv.ParseFloat(args[3], 64)

	p := &models.Payment{
		WarehouseID: wID,
		DistrictID: dID,
		CustomerID: cID,
		Amount: paymentAmt,
	}

	_, err := pc.s.ProcessTransaction(p)
	if err != nil {
		log.Println("error occurred in executing the payment transaction. Err: %v", err)
		return false
	}
	
	return true
}

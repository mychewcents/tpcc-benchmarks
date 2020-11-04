package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
)

type PaymentController interface {
	handler.TransactionHandler
}

type paymentControllerImpl struct {
	s service.PaymentService
}

func NewPaymentController(cassandraSession *common.CassandraSession) PaymentController {
	return &paymentControllerImpl{
		s: service.NewPaymentService(cassandraSession),
	}
}

func (p *paymentControllerImpl) HandleTransaction(cmd []string) {
	request := makePaymentRequest(cmd)
	p.s.ProcessPaymentTransaction(request)
	//printPaymentResponse(response)
}

func makePaymentRequest(cmd []string) *model.PaymentRequest {
	cWId, _ := strconv.Atoi(cmd[1])
	cDId, _ := strconv.Atoi(cmd[2])
	cId, _ := strconv.Atoi(cmd[3])
	payment, _ := strconv.ParseFloat(cmd[4], 64)

	return &model.PaymentRequest{
		CWId:    cWId,
		CDId:    cDId,
		CId:     cId,
		Payment: payment,
	}
}

func printPaymentResponse(r *model.PaymentResponse) {
	fmt.Println("*********************** Payment Transaction Output ***********************")
	fmt.Printf("1. Customer's identifier - (CWId: %v CDId: %v CId: %v), CName:%+v, CAddress:%+v CPhone:%v CSince:%v CCredit:%v CCreditLim:%.2f CDiscount:%.2f CBalance:%.2f\n", r.CWId, r.CDId, r.CId, r.CName, r.CAddress, r.CPhone, r.CSince, r.CCredit, r.CCreditLim, r.CDiscount, r.CBalance)
	fmt.Printf("2. Warehouse's address: %+v\n", r.WAddress)
	fmt.Printf("3. District's address: %+v\n", r.DAddress)
	fmt.Printf("4. Payment amount: %v\n", r.Payment)
	fmt.Println()
}

func (p *paymentControllerImpl) Close() error {
	panic("implement me")
}

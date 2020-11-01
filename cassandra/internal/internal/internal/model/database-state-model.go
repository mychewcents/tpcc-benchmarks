package model

import "github.com/gocql/gocql"

type DatabaseStateResponse struct {
	SumWYTD         float64
	SumDYTD         float64
	SumCBalance     float64
	SumCYTDPayment  float64
	SumCPaymentCnt  int64
	SumCDeliveryCnt int64
	MaxOId          gocql.UUID
	SumOOlCnt       int64
	SumOlAmount     float32
	SumOlQuantity   int64
	SumSQuantity    int64
	SumSYTD         int64
	SumSOrderCnt    int64
	SumSRemoteCnt   int64
}

package dbstate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// RecordDBState records the DB State for the experiment performed
func RecordDBState(db *sql.DB, experiment int, path string) error {
	var sumYTDWarehouse float64
	var sumYTDDistrict float64
	var sumNextOrderIDs int
	var cBalance, cYTDPayment float64
	var cPaymentCount, cDeliveryCount int
	var maxOrderID, totalOrderLineCount int
	var totalOrderAmount float64
	var totalOrderLineQuantity int
	var totalStockQuantity, totalYTD, totalOrderCount, totalRemoteCount int
	var err error

	if sumYTDWarehouse, err = getWarehouseState(db); err != nil {
		return fmt.Errorf("Error occured, getWarehouseState, Err: %v", err)
	}
	log.Printf("Warehouse State: %f ", sumYTDWarehouse)
	sumYTDDistrict = sumYTDWarehouse

	if sumNextOrderIDs, err = getDistrictState(db); err != nil {
		return fmt.Errorf("Error occured, getDistrictState, Err: %v", err)
	}
	log.Printf("District State: %f, %d ", sumYTDDistrict, sumNextOrderIDs)

	if cBalance, cYTDPayment, cPaymentCount, cDeliveryCount, err = getCustomerState(db); err != nil {
		return fmt.Errorf("Error occured, getCustomerState, Err: %v", err)
	}
	log.Printf("Customer State: %f, %f, %d, %d ", cBalance, cYTDPayment, cPaymentCount, cDeliveryCount)

	if maxOrderID, totalOrderLineCount, totalOrderAmount, err = getOrderState(db); err != nil {
		return fmt.Errorf("Error occured, getOrderState, Err: %v", err)
	}
	log.Printf("Order State: %d, %d, %f ", maxOrderID, totalOrderLineCount, totalOrderAmount)

	if totalOrderLineQuantity, err = getOrderLineState(db); err != nil {
		return fmt.Errorf("Error occured, getOrderLineState, Err: %v", err)
	}
	log.Printf("Order Line State: %d ", totalOrderLineQuantity)

	if totalStockQuantity, totalYTD, totalOrderCount, totalRemoteCount, err = getStockState(db); err != nil {
		return fmt.Errorf("Error occured, getOrderLineState, Err: %v", err)
	}
	log.Printf("Stock State: %d, %d, %d, %d ", totalStockQuantity, totalYTD, totalOrderCount, totalRemoteCount)

	outputStr := fmt.Sprintf("%d,%f,%f,%d,%f,%f,%d,%d,%d,%d,%f,%d,%d,%d,%d,%d",
		experiment,
		sumYTDWarehouse,
		sumYTDDistrict,
		sumNextOrderIDs,
		cBalance,
		cYTDPayment,
		cPaymentCount,
		cDeliveryCount,
		maxOrderID,
		totalOrderLineCount,
		totalOrderAmount,
		totalOrderLineQuantity,
		totalStockQuantity,
		totalYTD,
		totalOrderCount,
		totalRemoteCount,
	)

	log.Println(outputStr)
	csvFile, err := os.Create(fmt.Sprintf("%s/%d.csv", path, experiment))
	if err != nil {
		return fmt.Errorf("Error in creating CSV file, Err: %v", err)
	}
	defer csvFile.Close()

	if _, err := csvFile.WriteString(outputStr); err != nil {
		return fmt.Errorf("Error in writing the db state CSV file, Err: %v", err)
	}

	return nil
}

func getWarehouseState(db *sql.DB) (float64, error) {
	var sumYTD float64

	sqlStatement := `SELECT sum(D_YTD) FROM District`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&sumYTD); err != nil {
		return 0.0, err
	}
	return sumYTD, nil
}

func getDistrictState(db *sql.DB) (int, error) {
	var sumNextOrderIDs int

	sqlStatement := `SELECT sum(D_NEXT_O_ID) FROM District`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&sumNextOrderIDs); err != nil {
		return 0, err
	}
	return sumNextOrderIDs, nil
}

func getCustomerState(db *sql.DB) (float64, float64, int, int, error) {
	var cBalance, cYTDPayment float64
	var cPaymentCount, cDeliveryCount int

	sqlStatement := `SELECT sum(C_BALANCE), sum(C_YTD_PAYMENT), sum(C_PAYMENT_CNT), sum(C_DELIVERY_CNT) FROM Customer`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&cBalance, &cYTDPayment, &cPaymentCount, &cDeliveryCount); err != nil {
		return 0.0, 0.0, 0, 0, err
	}
	return cBalance, cYTDPayment, cPaymentCount, cDeliveryCount, nil
}

func getOrderState(db *sql.DB) (int, int, float64, error) {
	var tempMaxOrderID, maxOrderID, tempTotalOrderLineCount, totalOrderLineCount int
	var tempTotalOrderAmount, totalOrderAmount float64

	sqlStatement := `SELECT max(O_ID), sum(O_OL_CNT), sum(O_TOTAL_AMOUNT) FROM Orders_WID_DID`

	for w := 1; w < 11; w++ {
		for d := 1; d < 11; d++ {
			finalSQLStatement := strings.ReplaceAll(sqlStatement, "WID", strconv.Itoa(w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

			row := db.QueryRow(sqlStatement)
			if err := row.Scan(&tempMaxOrderID, &tempTotalOrderLineCount, &tempTotalOrderAmount); err != nil {
				return 0, 0, 0.0, err
			}

			if tempMaxOrderID > maxOrderID {
				maxOrderID = tempMaxOrderID
			}
			totalOrderLineCount += tempTotalOrderLineCount
			totalOrderAmount += tempTotalOrderAmount
		}
	}

	return maxOrderID, totalOrderLineCount, totalOrderAmount, nil
}

func getOrderLineState(db *sql.DB) (int, error) {
	var tempTotalQuantity, totalQuantity int

	sqlStatement := `SELECT sum(OL_QUANTITY) FROM Order_Line_WID_DID`

	for w := 1; w < 11; w++ {
		for d := 1; d < 11; d++ {
			finalSQLStatement := strings.ReplaceAll(sqlStatement, "WID", strconv.Itoa(w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

			row := db.QueryRow(sqlStatement)
			if err := row.Scan(&tempTotalQuantity); err != nil {
				return 0, err
			}

			totalQuantity += tempTotalQuantity
		}
	}
	return totalQuantity, nil
}

func getStockState(db *sql.DB) (int, int, int, int, error) {
	var totalQuantity, totalYTD, totalOrderCount, totalRemoteCount int

	sqlStatement := `SELECT sum(S_QUANTITY), sum(S_YTD), sum(S_ORDER_CNT), sum(S_REMOTE_CNT) FROM Stock`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&totalQuantity, &totalYTD, &totalOrderCount, &totalRemoteCount); err != nil {
		return 0, 0, 0, 0, err
	}

	return totalQuantity, totalYTD, totalOrderCount, totalRemoteCount, nil
}

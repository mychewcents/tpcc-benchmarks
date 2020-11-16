package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

//RelatedCustomerService interface to the related customer transaction
type RelatedCustomerService interface {
	ProcessTransaction(req *models.RelatedCustomer) (*models.RelatedCustomerOutput, error)
	Print(result *models.RelatedCustomerOutput)
}

type relatedCustomerServiceImpl struct {
	db *sql.DB
	ci dao.CustomerItemsPairDao
}

// CreateRelatedCustomerService creates new service
func CreateRelatedCustomerService(db *sql.DB) RelatedCustomerService {
	return &relatedCustomerServiceImpl{
		db: db,
		ci: dao.CreateCustomerItemsPairDao(db),
	}
}

// ProcessTransaction executes related customer transaction
func (rcs *relatedCustomerServiceImpl) ProcessTransaction(req *models.RelatedCustomer) (result *models.RelatedCustomerOutput, err error) {

	log.Printf("Starting the Related Customer Transaction for: w=%d d=%d c=%d", req.WarehouseID, req.DistrictID, req.CustomerID)

	result, err = rcs.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred in processing the related customer transaction. Err: %v", err)
	}

	log.Printf("Completed the Related Customer Transaction for: w=%d d=%d c=%d", req.WarehouseID, req.DistrictID, req.CustomerID)
	return
}

func (rcs *relatedCustomerServiceImpl) execute(req *models.RelatedCustomer) (result *models.RelatedCustomerOutput, err error) {
	itemPairsString, err := rcs.ci.GetItemPairsString(req.WarehouseID, req.DistrictID, req.CustomerID)
	if err != nil {
		return nil, err
	}

	for wID := 1; wID <= 10; wID++ {
		if req.WarehouseID != wID {
			result.Customers[wID] = make(map[int][]int)

			for dID := 1; dID <= 10; dID++ {
				customerIDs, err := rcs.ci.GetCustomerIDsWithSamePairs(req.WarehouseID, req.DistrictID, itemPairsString)
				if err != nil {
					return nil, err
				}
				if len(customerIDs) > 0 {
					result.Customers[wID][dID] = customerIDs
				}
			}
		}
	}

	return
}

func (rcs *relatedCustomerServiceImpl) Print(result *models.RelatedCustomerOutput) {
	var relatedCustomerOutputBuilder strings.Builder

	for wID, value := range result.Customers {
		for dID, customerIDs := range value {
			for _, value := range customerIDs {
				relatedCustomerOutputBuilder.WriteString(fmt.Sprintf("(%d, %d, %d),", wID, dID, value))
			}
		}
	}

	relatedCustomerOutputString := relatedCustomerOutputBuilder.String()
	relatedCustomerOutputString = relatedCustomerOutputString[:len(relatedCustomerOutputString)-1]

	fmt.Println(relatedCustomerOutputString)
}

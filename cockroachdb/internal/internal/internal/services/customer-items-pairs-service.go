package services

import (
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
)

// CustomerItemsPairService interface to initially create the customer items pairs service
type CustomerItemsPairService interface {
	LoadInitial(warehouses, districts int) error
}

type customerItemsPairServiceImple struct {
	db *sql.DB
	ci dao.CustomerItemsPairDao
	o  dao.OrderDao
	ol dao.OrderLineDao
}

// CreateCustomerItemsPairService creates a service to initially load the customer items pair table
func CreateCustomerItemsPairService(db *sql.DB) CustomerItemsPairService {
	return &customerItemsPairServiceImple{
		db: db,
		ci: dao.CreateCustomerItemsPairDao(db),
		o:  dao.CreateOrderDao(db),
		ol: dao.CreateOrderLineDao(db),
	}
}

func (cips *customerItemsPairServiceImple) LoadInitial(warehouses, districts int) (err error) {

	for w := 1; w <= warehouses; w++ {
		for d := 1; d <= districts; d++ {
			if err = cips.execute(w, d); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cips *customerItemsPairServiceImple) execute(wID, dID int) (err error) {

	customerOrders, err := cips.o.GetOrderIDForCustomers(wID, dID)
	if err != nil {
		return err
	}

	orderItems, err := cips.ol.GetDistinctItemIDsPerOrder(wID, dID)
	if err != nil {
		return err
	}

	err = cips.ci.InsertInitialPairs(wID, dID, customerOrders, orderItems)
	if err != nil {
		return err
	}

	return nil
}

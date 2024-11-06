package dmap

import (
	"fmt"
	"gRPCOrderService/pkg/order"
	"sync"
)

const NOT_FOUND = "order not found"

type DBOrder struct {
	orders map[string]*order.Order
	mux    sync.RWMutex
}

func New() *DBOrder {
	return &DBOrder{
		orders: make(map[string]*order.Order),
	} // TODO: Implement this function
}

func (db *DBOrder) AddOrder(oid string, order *order.Order) error {
	db.mux.Lock()
	db.orders[oid] = order
	db.mux.Unlock()

	return nil
}

func (db *DBOrder) GetOrder(oid string) (*order.Order, error) {
	db.mux.RLock()
	ord, ok := db.orders[oid]
	db.mux.RUnlock()

	if ok {
		return ord, nil
	} else {
		return nil, fmt.Errorf(NOT_FOUND)
	}
}

func (db *DBOrder) UpdateOrder(oid string, order *order.Order) error {
	if !db.ExistOrder(oid) {
		return fmt.Errorf(NOT_FOUND)
	}

	db.mux.Lock()
	db.orders[oid] = order
	db.mux.Unlock()

	return nil
}

func (db *DBOrder) DeleteOrder(oid string) error {
	if !db.ExistOrder(oid) {
		return fmt.Errorf(NOT_FOUND)
	}

	db.mux.Lock()
	delete(db.orders, oid)
	db.mux.Unlock()

	return nil
}

func (db *DBOrder) ExistOrder(oid string) bool {
	db.mux.RLock()
	defer db.mux.RUnlock()
	_, ok := db.orders[oid]
	return ok
}

func (db *DBOrder) RangeOrders(callback func(id string, order *order.Order)) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	for id, ord := range db.orders {
		callback(id, ord)
	}
}

package service

import (
	"fmt"
	"order/internal/types"
	"sync"
)

type OrderSvc struct {
	orderMap map[string]types.Order
	sync.RWMutex
}

func NewOrderSvc() OrderSvc {
	om := make(map[string]types.Order)
	return OrderSvc{orderMap: om}
}

func (o *OrderSvc) PlaceOrder(order types.Order) {
	o.Lock()
	o.orderMap[order.ID] = order
	o.Unlock()
}

func (o *OrderSvc) GetOrder(id string) (types.Order, error) {
	o.RLock()
	value, ok := o.orderMap[id]
	o.RUnlock()
	if ok {
		orderFromMap := value
		return orderFromMap, nil
	} else {
		no := types.Order{}
		return no, fmt.Errorf("Order not found")
	}
}

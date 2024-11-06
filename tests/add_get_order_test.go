package tests

import (
	"fmt"
	test "gRPCOrderService/pkg/api/3_1"
	order2 "gRPCOrderService/pkg/order"
	"gRPCOrderService/tests/suite"
	"testing"
)

const (
	appID     = 1             // ID приложения, которое мы создали миграцией
	appSecret = "test-secret" // Секретный ключ приложения
)

var order = &order2.Order{
	Item:     "Сочная, ароматная, горячая пицца",
	Quantity: 2,
}

func MatchOrders(fItem, sItem string, fQuantity, sQuantity int32) error {
	if fItem != sItem {
		return fmt.Errorf("expected item %s, got %s", fItem, sItem)
	} else if fQuantity != sQuantity {
		return fmt.Errorf("expected quantity %d, got %d", fQuantity, sQuantity)
	}
	return nil
}

func TestGetOrder(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	crOrder, err := st.Client.CreateOrder(ctx, &test.CreateOrderRequest{
		Item: order.Item, Quantity: order.Quantity}) // Создаём заказ

	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	t.Log("Order created")

	getOrder, err := st.Client.GetOrder(ctx, &test.GetOrderRequest{
		Id: crOrder.Id, // Идентификатор заказа
	})

	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	t.Log("Order got")

	if MatchOrders(order.Item, getOrder.Order.Item, order.Quantity, getOrder.Order.Quantity) != nil {
		t.Fatalf("orders do not match: %v", err)
	}
}

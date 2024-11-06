package tests

import (
	"sync"
	"testing"

	test "gRPCOrderService/pkg/api/3_1"
	"gRPCOrderService/tests/suite"
)

func TestDeleteOrderSuccess(t *testing.T) {
	ctx, st := suite.New(t)

	// Create an order to delete
	crOrder, err := st.Client.CreateOrder(ctx, &test.CreateOrderRequest{
		Item: "Test Item", Quantity: 1})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Attempt to delete the order
	delResp, err := st.Client.DeleteOrder(ctx, &test.DeleteOrderRequest{Id: crOrder.Id})
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}

	if !delResp.Success {
		t.Fatalf("expected success, got failure")
	}
}

func TestDeleteOrderNonExistentID(t *testing.T) {
	ctx, st := suite.New(t)

	// Attempt to delete a non-existent order
	_, err := st.Client.DeleteOrder(ctx, &test.DeleteOrderRequest{Id: "non-existent-id"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDeleteOrderConcurrentRequests(t *testing.T) {
	ctx, st := suite.New(t)

	// Create an order to delete
	crOrder, err := st.Client.CreateOrder(ctx, &test.CreateOrderRequest{
		Item: "Test Item", Quantity: 1})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	var wg sync.WaitGroup
	const numConcurrentRequests = 10

	for i := 0; i < numConcurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.Client.DeleteOrder(ctx, &test.DeleteOrderRequest{Id: crOrder.Id})
		}()
	}

	wg.Wait()

	// Verify the order is deleted
	_, err = st.Client.GetOrder(ctx, &test.GetOrderRequest{Id: crOrder.Id})
	if err == nil {
		t.Fatalf("expected error when getting deleted order, got nil")
	}
}

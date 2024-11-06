package grpc

import (
	"context"
	. "gRPCOrderService/internal/db"
	test "gRPCOrderService/pkg/api/3_1"
	"gRPCOrderService/pkg/order"
)

type OrderServiceServerAPI struct {
	test.UnimplementedOrderServiceServer
}

type OrderServiceAPI interface {
	CreateOrder(ctx context.Context, item string, number int) (string, error)
	GetOrder(ctx context.Context, id string) (*order.Order, error)
	UpdateOrder(ctx context.Context, id string, ord *order.Order) error
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) (map[string]*order.Order, error)
}

func NewOrderServiceServerAPI() *OrderServiceServerAPI {
	return &OrderServiceServerAPI{}
}

func (s *OrderServiceServerAPI) CreateOrder(ctx context.Context, in *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	id := GenerateID(in.Item, int(in.Quantity))
	if err := Db.AddOrder(id, &order.Order{Item: in.Item, Quantity: in.Quantity}); err != nil {
		return nil, err
	}
	return &test.CreateOrderResponse{Id: id}, nil
}

func (s *OrderServiceServerAPI) GetOrder(ctx context.Context, in *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	ord, err := Db.GetOrder(in.Id)
	if err != nil {
		return nil, err
	}
	return &test.GetOrderResponse{Order: &test.Order{
		Id:       in.Id,
		Item:     ord.Item,
		Quantity: ord.Quantity,
	}}, nil
}

func (s *OrderServiceServerAPI) UpdateOrder(_ context.Context, in *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	if err := Db.UpdateOrder(in.Id, &order.Order{
		Item:     in.Item,
		Quantity: in.Quantity,
	}); err != nil {
		return nil, err
	}

	return &test.UpdateOrderResponse{Order: &test.Order{
		Id:       in.Id,
		Item:     in.Item,
		Quantity: in.Quantity,
	}}, nil
}

func (s *OrderServiceServerAPI) DeleteOrder(_ context.Context, in *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	if err := Db.DeleteOrder(in.Id); err != nil {
		return nil, err
	}
	return &test.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderServiceServerAPI) ListOrders(context.Context, *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	var orders []*test.Order
	Db.RangeOrders(func(oid string, ord *order.Order) {
		orders = append(orders, &test.Order{
			Id:       oid,
			Item:     ord.Item,
			Quantity: ord.Quantity,
		})
	})

	return &test.ListOrdersResponse{Orders: orders}, nil
}

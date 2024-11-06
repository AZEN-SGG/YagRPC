package suite

import (
	"context"
	"gRPCOrderService/internal/config"
	grpc2 "gRPCOrderService/internal/transport/grpc"
	test "gRPCOrderService/pkg/api/3_1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"sync"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client test.OrderServiceClient
	Server *grpc.Server
}

var (
	once       sync.Once
	testSuite  *Suite
	serverAddr string
)

// startTestGRPCServer запускает gRPC-сервер для тестов один раз
func startTestGRPCServer(ctx context.Context, cfg *config.Config) (*grpc2.Server, error) {
	gRPCServer, err := grpc2.New(ctx, cfg.GRPCServerPort)
	if err != nil {
		return nil, err
	}

	go func() {
		if err = gRPCServer.Start(ctx); err != nil {
			panic(err)
		}
	}()

	return gRPCServer, nil
}

// InitializeTestSuite инициализирует тестовый набор один раз
func InitializeTestSuite(t *testing.T) *Suite {
	once.Do(func() {
		cfg := config.New("../configs/local.env")

		ctx, cancelCtx := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancelCtx()
		})

		// Запуск gRPC-сервера
		server, err := startTestGRPCServer(ctx, cfg)
		if err != nil {
			t.Fatalf("failed to start test gRPC server: %v", err)
		}

		t.Cleanup(func() {
			server.GrpcServer.Stop()
		})

		serverAddr = net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPCServerPort))

		// Создание клиента для gRPC
		cc, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("grpc server connection failed: %v", err)
		}

		client := test.NewOrderServiceClient(cc)

		// Создаем тестовый набор
		testSuite = &Suite{
			T:      t,
			Cfg:    cfg,
			Client: client,
			Server: server.GrpcServer,
		}
	})

	return testSuite
}

// New создает новый тестовый набор для каждого теста, используя уже инициализированный сервер и клиент
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// Инициализируем тестовый набор только один раз
	suite := InitializeTestSuite(t)

	// Создаем отдельный контекст для каждого теста
	ctx, cancelCtx := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelCtx()
	})

	return ctx, suite
}

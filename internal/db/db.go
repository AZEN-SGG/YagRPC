package db

import (
	"fmt"
	"gRPCOrderService/pkg/db/dmap"
	"time"
)

var Db = dmap.New()

func GenerateID(dishName string, quantity int) string {
	timestamp := time.Now().Format("20060102150405.000")
	id := fmt.Sprintf("%s.%d.%s", dishName, quantity, timestamp)

	return id
}

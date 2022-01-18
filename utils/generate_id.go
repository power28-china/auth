package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/power28-china/auth/utils/logger"
)

var (
	node *snowflake.Node
	err  error
)

func init() {
	node, err = snowflake.NewNode(1)
	if err != nil {
		logger.Sugar.Fatalf("err")
		panic(err)
	}
}

// GetID return a snowflake ID.
func GetID() int64 {
	return node.Generate().Int64()
}

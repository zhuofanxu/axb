package uniqueidutils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	// config time: approximately at 2024-08-11 11:40:00 (UTC+8)
	// Set the epoch to 1989-10-08 15:52:18 (UTC+8), ensuring the elapsed time starts near the minimum value of 41 bits (1099511627776).
	// With a 41-bit timestamp, the maximum value of elapsed time (2199023255551) will be reached after approximately 35 years,
	// corresponding to the date 2059-06-15 07:39:53 (UTC+8).
	snowflake.Epoch = 623836338208
	node, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GenSnowflakeId() snowflake.ID {
	return node.Generate()
}

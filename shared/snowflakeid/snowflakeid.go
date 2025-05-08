package snowflakeid

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

type SnowflakeInterface interface {
	Generate() (r int64)
}

type SnowflakeImpl struct {
	snowflakeID *snowflake.Node
}

func NewSnowflake(snowflakeID *snowflake.Node) *SnowflakeImpl {
	return &SnowflakeImpl{
		snowflakeID,
	}
}

func CreateSnowflakeNode(nodeNumber int64) *snowflake.Node {
	node, err := snowflake.NewNode(nodeNumber)
	if err != nil {
		panic(fmt.Sprintf("snowflake error: %v\n", err))
	}
	return node
}

var _ SnowflakeInterface = (*SnowflakeImpl)(nil)

func (s *SnowflakeImpl) Generate() (r int64) {
	return s.snowflakeID.Generate().Int64()
}

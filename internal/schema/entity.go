package schema

import (
	"strconv"
)

type EntityId int64

func (i EntityId) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i EntityId) Int64() int64 {
	return int64(i)
}

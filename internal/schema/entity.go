package schema

import (
	"strconv"
)

type EntityId int64

func (i EntityId) String() string {
	return strconv.FormatInt(int64(i), 10)
}

package schema

import (
	"strconv"
)

type EntityId uint

func (i EntityId) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

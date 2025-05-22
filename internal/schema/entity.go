package schema

import (
	"strconv"
)

type EntityId uint

//type EntityId int64

func (i EntityId) String() string {
	return strconv.FormatUint(uint64(i), 10)
	//return strconv.FormatInt(int64(i), 10)
}
func (i EntityId) Uint() uint {
	return uint(i)
}

// EntityIdFromString 从字符串转换为EntityId
func EntityIdFromString(s string) (EntityId, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return EntityId(id), nil
}

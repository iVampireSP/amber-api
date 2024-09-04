package schema

import (
	"database/sql/driver"
	"github.com/bytedance/sonic"
)

type Embedding []float32

func (e *Embedding) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), e)
}

func (e Embedding) Value() (driver.Value, error) {
	return sonic.Marshal(e)
}

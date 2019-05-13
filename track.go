package god

import (
	"fmt"
	"reflect"
)

type tracker struct {
	value reflect.Value
}

// Track will track the length of a value.
func Track(name string, value interface{}) {
	add(name, &tracker{
		value: reflect.ValueOf(value),
	})
}

func (c *tracker) string() string {
	return fmt.Sprintf("%d", c.value.Len())
}

func (c *tracker) reset() {}

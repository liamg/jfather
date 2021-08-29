package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeArray(v reflect.Value) error {

	switch v.Kind() {
	case reflect.Array:
	case reflect.Slice:
		out.Set(reflect.MakeSlice(out.Type(), l, l))
	}

	return fmt.Errorf("not implemented")
}

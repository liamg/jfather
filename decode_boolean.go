package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeBoolean(v reflect.Value) error {
	if v.Kind() != reflect.Bool {
		return fmt.Errorf("cannot decode boolean value to non-boolean target")
	}

	v.SetBool(n.raw.(bool))
	return nil
}

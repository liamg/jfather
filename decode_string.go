package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeString(v reflect.Value) error {

	if v.Kind() != reflect.String {
		return fmt.Errorf("cannot decode string value to non-string target: %s", v.Kind())
	}

	v.SetString(n.raw.(string))
	return nil
}

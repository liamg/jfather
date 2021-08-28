package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeBoolean(target interface{}) error {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot decode to non-pointer target")
	}

	if v.IsNil() {
		if !v.CanAddr() {
			return fmt.Errorf("cannot write to unaddressable value")
		}
		iv := n.raw.(bool)
		v.Set(reflect.ValueOf(&iv))
		return nil
	}

	if v.Elem().Kind() != reflect.Bool {
		return fmt.Errorf("cannot decode string value to non-string target")
	}

	v.Elem().SetBool(n.raw.(bool))
	return nil
}

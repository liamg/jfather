package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeString(target interface{}) error {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot decode to non-pointer target")
	}

	if v.IsNil() {
		if !v.CanAddr() {
			return fmt.Errorf("cannot write to unaddressable value")
		}
		iv := n.raw.(string)
		v.Set(reflect.ValueOf(&iv))
		return nil
	}

	if v.Elem().Kind() != reflect.String {
		return fmt.Errorf("cannot decode string value to non-string target")
	}

	v.Elem().SetString(n.raw.(string))
	return nil
}

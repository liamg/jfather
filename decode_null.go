package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeNull(target interface{}) error {

	v := reflect.ValueOf(target)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot decode to non-pointer target")
	}

	p := v.Elem()
	p.Set(reflect.Zero(p.Type()))
	return nil
}

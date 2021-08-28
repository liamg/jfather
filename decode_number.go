package jfather

import (
	"fmt"
	"reflect"
)

func (n *node) decodeNumber(target interface{}) error {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot decode to non-pointer target")
	}

	switch v.Elem().Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		if i64, ok := n.raw.(int64); ok {
			v.Elem().SetInt(i64)
			return nil
		}
		if f64, ok := n.raw.(float64); ok {
			v.Elem().SetInt(int64(f64))
			return nil
		}
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		if i64, ok := n.raw.(int64); ok {
			v.Elem().SetUint(uint64(i64))
			return nil
		}
		if f64, ok := n.raw.(float64); ok {
			v.Elem().SetUint(uint64(f64))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if i64, ok := n.raw.(int64); ok {
			v.Elem().SetFloat(float64(i64))
			return nil
		}
		if f64, ok := n.raw.(float64); ok {
			v.Elem().SetFloat(f64)
			return nil
		}
	default:
		return fmt.Errorf("cannot decode number value to *%s target", v.Elem().Kind())
	}

	return fmt.Errorf("internal value is not numeric")
}

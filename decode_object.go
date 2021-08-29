package jfather

import (
	"fmt"
	"reflect"
	"strings"
)

func (n *node) decodeObject(v reflect.Value) error {

	if v.Kind() == reflect.Struct {
		return n.decodeObjectToStruct(v)
	}

	if v.Kind() == reflect.Map {
		return fmt.Errorf("not supported")
	}

	return fmt.Errorf("bad target type")
}

func (n *node) objectAsMap() (map[string]Node, error) {
	if n.kind != KindObject {
		return nil, fmt.Errorf("not an object")
	}
	properties := make(map[string]Node)
	contents := n.content
	for i := 0; i < len(contents); i += 2 {
		key := contents[i]
		if key.Kind() != KindString {
			return nil, fmt.Errorf("invalid object key - please report this bug")
		}
		keyStr := key.(*node).raw.(string)

		if i+1 >= len(contents) {
			return nil, fmt.Errorf("missing object value - please report this bug")
		}
		properties[keyStr] = contents[i+1]
	}
	return properties, nil
}

func (n *node) decodeObjectToStruct(v reflect.Value) error {

	properties, err := n.objectAsMap()
	if err != nil {
		return err
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fv := t.Field(i)
		tags := strings.Split(fv.Tag.Get("json"), ",")
		var tagName string
		for _, tag := range tags {
			if tag != "omitempty" && tag != "-" {
				tagName = tag
			}
		}
		if tagName == "" {
			tagName = fv.Name
		}

		value, ok := properties[tagName]
		if !ok {
			// TODO: should we zero this value?
			continue
		}

		subject := v.Field(i)

		// if fields are nil pointers, initialise them with values of the correct type
		if subject.Kind() == reflect.Ptr && subject.IsNil() {
			subject.Set(reflect.New(subject.Type().Elem()))
		}

		if err := value.(*node).decodeToValue(subject); err != nil {
			return err
		}
	}
	return nil
}

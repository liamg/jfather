package main

type Unmarshaller interface {
	UnmarshalJFather(node Node) error
}

func Unmarshal(data []byte, target interface{}) error {
	node := parse(data)
	if unmarshaller, ok := target.(Unmarshaller); ok {
		return unmarshaller.UnmarshalJFather(node)
	}
	return node.Decode(target)
}

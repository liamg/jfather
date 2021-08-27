package main

func parse(data []byte) Node {
	n := &node{
		data: data,
	}
	return n
}

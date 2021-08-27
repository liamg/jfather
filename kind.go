package main

type Kind uint8

const (
	KindUnknown Kind = iota
	KindNumber
	KindString
	KindList
	KindObject
	KindNull
)

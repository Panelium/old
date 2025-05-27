package main

import proto_gen_go "panelium/proto-gen-go"

func main() {
	test := proto_gen_go.TestMessage{}

	text := "Hello, World!"
	test.Text = &text

	number := int32(42)
	test.Number = &number

	array := []string{"item1", "item2", "item3"}
	test.Array = array

	boolean := true
	test.Boolean = &boolean

	println(test.String())
}

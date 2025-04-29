package main

import (
	"fmt"
)

func main() {
	fmt.Println("test")
	fmt.Println("test", "but with additional argument")
	fmt.Println(42, "number and string")
	fmt.Println("nothing here")
	first_function := func() {
		fmt.Println("hello from function")
		fmt.Println("hello from function")

	}
	first_function()
	fmt.Println(10 + 12)
	fmt.Println(10 - 12)
	fmt.Println(10 * 12)
	fmt.Println(12 / 6)
}

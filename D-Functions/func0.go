// Variable type comes after the variable name string = Joseph
package main

import "fmt"

func sum(x int, y int) int {
	return x + y
}

// Define the last arg the type
func sub(x, y int) int {
	return x - y
}

func concat0(sent1 string, sent2 string) string {
	return sent1 + sent2
}

func main() {
	fmt.Println(sum(5, 2))
	fmt.Println(sub(5, 2))
	fmt.Println(concat0("Hi Yuksel,", " Received 5 job offers"))
}

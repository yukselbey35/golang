package main

import "fmt"

func main() {
	age := 7
	fmt.Printf("You are %d years old\n", age)
	greeting := "Hello"
	num1 := 39
	fmt.Printf("%s your age was %v when you came\n", greeting, num1)

	const name = "Joseph"
	const rate_S = 99.5

	message := fmt.Sprintf("Hi %s, your app success is %.2f percent", name, rate_S)

	// don't edit below this line

	fmt.Println(message)
}

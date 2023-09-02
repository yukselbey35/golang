package main

import "fmt"

func main() {
	email := "Hi, How are you doing today"
	length := len(email)
	if length < 1 {
		fmt.Println("Email is invalid")
	} else {
		fmt.Println("Email is valid")
	}
}

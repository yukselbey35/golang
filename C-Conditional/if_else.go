package main

import "fmt"

func main() {
	temp := 90
	if temp <= 60 {
		fmt.Println("Cold day")
	} else if temp <= 85 {
		fmt.Println("Weather is cool.")
	} else {
		fmt.Println("Warm day")
	}
}

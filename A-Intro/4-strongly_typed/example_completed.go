// we can't concatenated with + because string and int, convert to int to string
package main

import "fmt"

func main() {
	var username string = "yuksel"
	var password string = "123456"

	fmt.Println("Authentication: ", username+": "+password)
}

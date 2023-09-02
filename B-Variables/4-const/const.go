// Use const when need unchangeable and read-only values
package main

import "fmt"

func main() {
	const name_Node = "node1"
	const CPU_thresHold = 80
	fmt.Println(name_Node, ":", CPU_thresHold, "%")
}

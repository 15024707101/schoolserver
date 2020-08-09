package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Remove("static/img/headPhoto/2cd0cb2f-1569-4e0f-aff3-15c16d7b7b97.jpg")
	if err != nil {
		fmt.Printf("remove ./file1.txt err : %v\n", err)
	}

}

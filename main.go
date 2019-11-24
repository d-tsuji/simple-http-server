package main

import (
	"fmt"
)

func main() {
	for {
		if err := Run(); err != nil {
			fmt.Printf("%+v", err)
		}
	}
}

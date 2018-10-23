package main

import (
	"fmt"
)

func main() {
	nc := collectContributions()
	fmt.Printf("\nCollected %v contributions\n\n", len(nc))
	for _, c := range nc {
		fmt.Printf("Contribution from %v into %v\n", c.User, c.Project)
	}
}

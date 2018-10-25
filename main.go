package main

import (
	"fmt"
)

func main() {
	nc := collectContributions()
	fmt.Printf("\nCollected %v contributions\n\n", len(nc))
	fmt.Print("\nSending Email...\n")
	sendMonthlyEmail(nc)
	fmt.Println("Email sent!")
}

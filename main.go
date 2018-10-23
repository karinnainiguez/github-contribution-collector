package main

import (
	"fmt"
)

func main() {
	nc := newClient()

	repos := getRepos(nc, "kubernetes-sigs")
	fmt.Println(len(repos))
}

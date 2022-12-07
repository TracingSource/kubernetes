package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Effective uid: %d\n", os.Geteuid())
}

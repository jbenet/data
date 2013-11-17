package main

import (
	"fmt"
)

const VERSION = "0.0.0"

func VersionCmd() {
	fmt.Println("data version", VERSION)
}

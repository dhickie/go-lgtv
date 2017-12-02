package main

import (
	"fmt"

	"github.com/dhickie/go-lgtv/discovery"
)

func main() {
	ip, _ := discovery.Discover("192.168.1.0")

	fmt.Printf("%v", ip)
}

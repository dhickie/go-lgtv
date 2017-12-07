package main

import (
	"fmt"
	"net"

	"github.com/dhickie/go-lgtv/control"
)

func main() {
	// tv, err := discovery.Discover("192.168.1.1")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	tv := control.NewTV(net.IP{192, 168, 1, 129})
	_, err := tv.Connect("45d55ed2e385752d6f6a86178d50a682")
	if err != nil {
		fmt.Println(err)
	}
	err = tv.GetChannelList()
	if err != nil {
		fmt.Println(err)
	}
}

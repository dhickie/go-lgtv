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
	_, err := tv.Connect("7668cb15d16a1a319f3731a9264b700b")
	if err != nil {
		fmt.Println(err)
	}
	err = tv.SwitchInput("HDMI_2")
	if err != nil {
		fmt.Println(err)
	}
}

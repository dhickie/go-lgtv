package control

import "net"

// LgTv represents the TV being controlled
type LgTv struct {
	ip  net.IP
	pin string
}

// NewTV returns a new LgTv object with the specified IP and pin
func NewTV(ip net.IP, pin string) LgTv {
	return LgTv{ip, pin}
}

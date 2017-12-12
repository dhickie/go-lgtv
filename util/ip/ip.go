package ip

import (
	"net"
	"strconv"
	"strings"
)

// ParseIP parses a string in to an IPv4 address
func ParseIP(ipStr string) (net.IP, error) {
	s := strings.Split(ipStr, ".")
	ip := net.IP{0, 0, 0, 0}

	for i := 0; i < 4; i++ {
		v, err := strconv.ParseInt(s[i], 10, 0)
		if err != nil {
			return nil, err
		}

		ip[i] = byte(v)
	}

	return ip, nil
}

package discovery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dhickie/go-lgtv/control"
	xmlutil "github.com/dhickie/go-lgtv/util/xml"
)

const openPort = 1042

var (
	// ErrNoTVFound indicates that no TV could be found
	ErrNoTVFound  = errors.New("Failed to find a TV")
	errNoResponse = errors.New("No response from IP address")
)

// Discover searches for LG TVs running version 3.5 of WebOS.
//
// This version has port 1042 open to GET requests to display details
// of the TV. gatewayIP is the IP address of the default gateway on the local network.
//
// This has only been tested using the C7V 2017 model.
func Discover(gatewayIP string) (control.LgTv, error) {
	// Convert the provided string in to an IP address
	gwIP, err := extractIP(gatewayIP)
	if err != nil {
		return control.LgTv{}, err
	}

	// Iterate over all possible local IP addresses (based on a single gateway setup)
	ips := make([]net.IP, 256)
	for i := range ips {
		ips[i] = net.IP{gwIP[0], gwIP[1], gwIP[2], byte(i)}
	}

	wg := sync.WaitGroup{}
	results := make([]bool, 256)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		startIndex := i * 32
		endIndex := (i + 1) * 32
		go pingWorker(ips[startIndex:endIndex], results[startIndex:endIndex], &wg)
	}
	wg.Wait()

	for i, v := range results {
		if v {
			gwIP[3] = byte(i)
			return control.NewTV(gwIP, ""), nil
		}
	}

	return control.LgTv{}, ErrNoTVFound
}

func pingWorker(ips []net.IP, results []bool, wg *sync.WaitGroup) {
	for i, v := range ips {
		result, err := pingIP(v)
		if err != nil || !result {
			results[i] = false
		} else {
			results[i] = true
		}
	}

	wg.Done()
}

func pingIP(ip net.IP) (bool, error) {
	timeout := time.Duration(20 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	resp, _ := client.Get(fmt.Sprintf("http://%v:%v", ip, openPort))
	if resp == nil {
		return false, errNoResponse
	}

	body, _ := ioutil.ReadAll(resp.Body)

	node, err := xmlutil.FindXMLNode(string(body), "modelName")
	if err != nil {
		if err != xmlutil.ErrNodeNotFound {
			return false, err
		}

		return false, nil
	}

	if string(node.Content) == "LG Smart TV" {
		return true, nil
	}
	return false, nil
}

func extractIP(ipStr string) (net.IP, error) {
	s := strings.Split(ipStr, ".")
	ip := net.IP{0, 0, 0, 0}

	for i := 0; i < 3; i++ {
		v, err := parseInt(s[i])
		if err != nil {
			return nil, err
		}

		ip[i] = byte(v)
	}

	return ip, nil
}

func parseInt(i string) (int64, error) {
	return strconv.ParseInt(i, 10, 0)
}

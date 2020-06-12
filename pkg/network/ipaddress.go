package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetMyIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("get my IP address failed")
}

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

func GetExternalIP() (string, error) {
	resp, err := http.Get("http://whatismyip.akamai.com/")
	if err != nil {
		resp, err = http.Get("http://ifconfig.io/ip")
		if err != nil {
			resp, err = http.Get("http://tnx.nl/ip")
			if err != nil {
				return "", err
			}
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

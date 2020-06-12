// Finding named hosts on a network

package main

import (
	"fmt"
	"os"
	"tools/pkg/network"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}
}

func main() {

	myip, err := network.GetMyIP()
	checkErr(err)
	fmt.Printf("LAN IP address: %s\n", myip)

	extip, err := network.GetExternalIP()
	checkErr(err)
	fmt.Printf("WAN IP address: %s\n", extip)

	machines := network.DiscoverLan(myip)
	fmt.Printf("machines on LAN\n")
	for _, machine := range machines {
		fmt.Printf("%s\n", machine)
	}

	ports := network.GetOpenPorts()
	fmt.Printf("listen ports on machine\n")
	for _, port := range ports {
		fmt.Printf("%s\n", port)
	}
}

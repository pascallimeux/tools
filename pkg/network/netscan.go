package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Machines struct {
	M   []string
	mux sync.Mutex
}

func (m *Machines) addMachine(IPAddress string) {
	m.M = append(m.M, IPAddress)
}

func DiscoverLan(IPAddress string) []string {
	subnetToScan := IPAddress[:strings.LastIndex(IPAddress, ".")]

	machines := Machines{}
	activeThreads := 0
	doneChannel := make(chan bool)

	for ip := 0; ip <= 255; ip++ {
		fullIP := subnetToScan + "." + strconv.Itoa(ip)
		go resolve(&machines, fullIP, doneChannel)
		activeThreads++
	}

	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
	return machines.M

}

func resolve(machines *Machines, ip string, doneChannel chan bool) {
	addresses, err := net.LookupAddr(ip)
	if err == nil {
		machines.mux.Lock()
		machines.addMachine(fmt.Sprintf("%s - %s", ip, strings.Join(addresses, ", ")))
		machines.mux.Unlock()
		//fmt.Printf("%s - %s\n", ip, strings.Join(addresses, ", "))
	}
	doneChannel <- true
}

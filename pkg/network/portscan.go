package network

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
}

type OpenPorts struct {
	P   []string
	mux sync.Mutex
}

func (op *OpenPorts) addPort(port string) {
	op.P = append(op.P, port)
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func ScanPort(op *OpenPorts, ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(op, ip, port, timeout)
		}
		return
	}

	conn.Close()
	op.mux.Lock()
	op.addPort(strconv.Itoa(port))
	//fmt.Println(port, "open")
	op.mux.Unlock()
}

func (ps *PortScanner) Start(op *OpenPorts, f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(op, ps.ip, port, timeout)
		}(port)
	}
}

func GetOpenPorts() []string {
	op := OpenPorts{}
	ps := &PortScanner{
		ip:   "127.0.0.1",
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(&op, 1, 65535, 500*time.Millisecond)
	return op.P
}

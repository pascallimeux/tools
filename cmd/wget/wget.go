package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	url         = "http://wwww.google.fr"
	mode        = int(0666)
	timeout     = 7 * time.Second
	output_path = "test.txt"
)

func init() {
	flag.StringVar(&url, "url", url, "")
	flag.DurationVar(&timeout, "timeout", timeout, "")
	flag.StringVar(&output_path, "output_path", output_path, "")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Printf("Usage: ./wget [-url] [-output] [-timeout]\n")
	fmt.Printf("   %-20.20s    %-70.70s\n", "-url string", fmt.Sprintf("url (default \"%s\")", url))
	fmt.Printf("   %-20.20s    %-70.70s\n", "-timeout duration", fmt.Sprintf("timeout (default %ds)", timeout/1000000000))
	fmt.Printf("   %-20.20s    %-70.70s\n", "-output string", fmt.Sprintf("output path (default \"%s\")", output_path))
}

func HTTPGet(url string, timeout time.Duration) (content []byte, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	ctx, cancel_func := context.WithTimeout(context.Background(), timeout)
	request = request.WithContext(ctx)

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		cancel_func()
		return nil, fmt.Errorf("INVALID RESPONSE; status: %s", response.Status)
	}

	return ioutil.ReadAll(response.Body)
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "agent" {
		url = fmt.Sprintf("http://192.168.20.53:8080/ccp3/%s/agent", runtime.GOOS)
		timeout = 10 * time.Second
		output_path = "agent"
		mode = int(0764)
	}
	fmt.Printf("download %s -> %s\n", url, output_path)
	content, err := HTTPGet(url, timeout)
	if err != nil {
		log.Fatalln("HTTPGET: ", err)
	}
	err = ioutil.WriteFile(output_path, content, os.FileMode(mode))
	if err != nil {
		log.Fatalln("WriteFile: ", err)
	}
	return
}

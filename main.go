package main

import (
	"flag"
	"fmt"
	"github.com/musicglue/doozer"
	"net"
	"strings"
	"time"
)

type empty struct{}

var (
	port        = flag.String("port", "8046", "The port to test the server on")
	servers     = flag.String("servers", "localhost,localhost,localhost", "The servers to test, either strings or IPs")
	protocol    = flag.String("protocol", "udp", "Use TCP or UDP")
	ownIp       = flag.String("ownip", "127.0.0.1", "The originating IP of the interface you are connecting from.")
	timeout     = flag.Int("timeout", 1000, "Timeout in Milliseconds before dialling the next server.")
	clustername = flag.String("cluster", "local", "Name of the cluster, local is default")
	output      string
)

func appendIfMissing(slice []string, s string) (output []string) {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func checkServer(name string) {
	fmt.Println(name)
	return
}

func main() {
	flag.Parse()

	serverList := strings.Split(*servers, ",")
	res := make([]string, 0, 0)
	sem := make(chan empty, len(serverList))

	for ind, server := range serverList {
		go func(ind int, name string) {
			address := name + ":" + *port
			_, err := net.DialTimeout(*protocol, address, time.Duration(*timeout)*time.Millisecond)
			if err != nil {
			} else {
				res = appendIfMissing(res, address)
			}
			sem <- empty{}
		}(ind, server)
	}

	for i := 0; i < len(serverList); i++ {
		<-sem
	}

	doozer, derr := doozer.Dial(res[0])
	if derr != nil {
		panic(derr)
		// HANDLE THIS PROPERLY
	}

	// Get latest Revision of DB
	rev, _ := doozer.Rev()
	// HANDLE THE ERROR

	// Get NS key to query
	key, _ := doozer.Getdir("/ctl/ns/"+*clustername, rev, 0, 1)
	// HANDLE THE ERROR

	query := fmt.Sprintf("/ctl/ns/%s/%s", *clustername, key[0])
	// Get IP Address to bind to
	value, _, _ := doozer.Get(query, nil)
	// HANDLE THE ERROR

	if value == *ownIp {
		output = res[0]
	} else {
		output = value
	}

	fmt.Printf("%s", output)

}

package main

import (
	"flag"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
)

const (
	endpointuri = "/api/v1/namespaces/default/endpoints/"
)

func main() {
	var (
		server    string
		//frequency int
	)

	/* Handling flags */
	flag.StringVar(&server, "server", "http://192.168.3.130:2181", "url for zk  server, e.g., http://192.168.3.130:2181")
	//flag.IntVar(&frequency, "frequency", 5, "watch frequency")
	flag.Parse()

	//ZK
	c, _, err := zk.Connect([]string{"192.168.3.130"}, time.Second) //*10)
	if err != nil {
		fmt.Printf("Error connecting to zk: %+v", err)
		//panic(err)
	}

	_, _, ch, err := c.ChildrenW("/ActiveStreams")
	if err != nil {
		fmt.Printf("== err: %+v \n", err)
	}

	fmt.Printf("==== Start Watching ZK ====! \n")
	for {
		select {
		case event := <- ch:
			fmt.Printf("==== [%v] - event: %+v \n", time.Now(), event)
		}
	}
}

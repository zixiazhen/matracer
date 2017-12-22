package main

import (
	"flag"
	"fmt"
	"time"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	endpointuri = "/api/v1/namespaces/default/endpoints/"
)

func main() {
	var (
		apiserver    string
		//frequency int
	)

	/* Handling flags */
	flag.StringVar(&apiserver, "apiserver", "http://127.0.0.1:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	//flag.IntVar(&frequency, "frequency", 5, "watch frequency")
	flag.Parse()


	//ZK
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	children, stat, ch, err := c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)




}

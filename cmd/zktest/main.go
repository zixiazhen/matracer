package main

import (
	"flag"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
)

func main() {
	var (
		server string
		//frequency int
	)

	/* Handling flags */
	flag.StringVar(&server, "server", "http://192.168.3.130:2181", "url for zk  server, e.g., http://192.168.3.130:2181")
	//flag.IntVar(&frequency, "frequency", 5, "watch frequency")
	flag.Parse()

	//ZK
	c, _, err := zk.Connect([]string{"192.168.3.130:2181"}, time.Second) //*10)
	if err != nil {
		fmt.Printf("Error connecting to zk: %+v", err)
		//panic(err)
	}

	state := c.State()
	fmt.Printf("state: %+v", state)


	createNode := zk.CreateRequest{Path: "/ActiveStreams/s01/DONE", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	createNode2 := zk.CreateRequest{Path: "/ActiveStreams/s01/node01", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	deleteNode := zk.DeleteRequest{Path: "/ActiveStreams/s01/DONE"}
	deleteWorkerRqsts := []interface{}{&createNode, &createNode2, &deleteNode}
	_, err = c.Multi(deleteWorkerRqsts...)
	if err != nil {
		fmt.Printf("====  delete error! %+v \n", err)
	}


	/*
	data, stat, err := c.Get("/")
	if  err != nil {
		fmt.Printf("Get returned error: %+v", err)
	}

	fmt.Printf("data: %+v \n", data)
	fmt.Printf("stat: %+v \n", stat)
	*/


/*	exists, stat, chl, err := c.ExistsW("/test/1")
	if !exists {
		fmt.Printf("%+v", "node is not there!! chl: %+v", chl)
		//return
	}
	for {
		select{
		case event := <- chl:
			fmt.Printf("event: %+v \n", event)
			time.Sleep(2 * time.Second)
		}
	}*/




/*	flags := int32(zk.FlagSequence | zk.FlagEphemeral)
	//flags := int32(-1)

	path := "/ActiveStreams/s01/DONE"
	p1, err := c.Create(
		path,
		[]byte{1, 2, 3, 4},
		flags,
		zk.WorldACL(zk.PermAll));
	if err != nil {
		fmt.Printf("==== Create  error: %+v \n", err)
	}
	fmt.Printf("==== Node created! %+v \n", p1)*/


	//deleteNode := zk.DeleteRequest{Path: "/ActiveStreams/s01/node01"}
	//createNode := zk.CreateRequest{Path: "/ActiveStreams/s01/node01"}
	//deleteWorkerRqsts := []interface{}{&createNode, &deleteNode}
	//createWorkerRqsts := []interface{}{&createNode}
	//_, err = c.Multi(createWorkerRqsts...)
	//if err != nil {
	//	fmt.Printf("====  create error! %+v \n", err)
	//}


/*
	createNode := zk.CreateRequest{Path: "/ActiveStreams/s01/DONE", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	createNode2 := zk.CreateRequest{Path: "/ActiveStreams/s01/node01", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	deleteNode := zk.DeleteRequest{Path: "/ActiveStreams/s01/DONE"}
	deleteWorkerRqsts := []interface{}{&createNode2, &createNode, &deleteNode}
	_, err = c.Multi(deleteWorkerRqsts...)
	if err != nil {
		fmt.Printf("====  delete error! %+v \n", err)
	}
*/

	//createNode1 := zk.CreateRequest{Path: "/ActiveStreams/s01/node01", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	//deleteNode := zk.DeleteRequest{Path: "/ActiveStreams/s01/node01"}
	//deleteWorkerRqsts := []interface{}{&createNode1, &deleteNode}
	//_, err = c.Multi(deleteWorkerRqsts...)
	//if err != nil {
	//	fmt.Printf("====  delete error! %+v \n", err)
	//}

	//createStreams(c)

	/*
	flags := int32(zk.FlagSequence | zk.FlagEphemeral)

	path := "/ActiveStreams/1/le_"
	p1, err := c.Create(
		path,
		[]byte{1, 2, 3, 4},
		flags,
		zk.WorldACL(zk.PermAll));
	if err != nil {
		fmt.Printf("==== Create  error: %+v \n", err)
	}
	fmt.Printf("==== Node created! %+v \n", p1)

	path = "/ActiveStreams/2/le_"
	p2, err := c.CreateProtectedEphemeralSequential(
		path,
		[]byte{1, 2, 3, 4},
		zk.WorldACL(zk.PermAll));
	if err != nil {
		fmt.Printf("==== Create  error: %+v \n", err)
	}
	fmt.Printf("==== Node created! %+v \n", p2)

	for {
		select {}
	}

	fmt.Printf("==== Done! \n")
	*/
}

/*

func createStreams(c *zk.Conn) {

	path0 := "/ActiveStreams"
	p0, _, _ := c.Get(path0)
	if len(p0) == 0 {
		c.Create(
			path0,
			[]byte{1, 2, 3, 4},
			0,
			zk.WorldACL(zk.PermAll));
	}

	path1 := "/ActiveStreams/1"
	p1, _, _ := c.Get(path1)
	if len(p1) == 0 {
		c.Create(
			path1,
			[]byte{1},
			0,
			zk.WorldACL(zk.PermAll));
	}

	path2 := "/ActiveStreams/2"
	p2, _, _ := c.Get(path2)
	if len(p2) == 0 {
		c.Create(
			path2,
			[]byte{2},
			0,
			zk.WorldACL(zk.PermAll));
	}

}
*/

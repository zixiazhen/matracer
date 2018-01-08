package main

import (
	"flag"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
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

	createStreams(c)

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
		select {
		}
	}

	fmt.Printf("==== Done! \n")
}

func createStreams(c *zk.Conn) {

	path0 := "/ActiveStreams"
	p0, _, _ := c.Get(path0)
	if len(p0) == 0 {
		c.Create(
			path0,
			[]byte{1,2,3,4},
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

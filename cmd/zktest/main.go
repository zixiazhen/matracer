package main

import (
	"flag"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"sync"
	//"reflect"
	//"sort"
)

func main() {
	var (
		server string
		//frequency int
	)

	/* Handling flags */
	flag.StringVar(&server, "server", "http://0.0.0.0:2181", "url for zk  server, e.g., http://192.168.3.130:2181")
	//flag.IntVar(&frequency, "frequency", 5, "watch frequency")
	flag.Parse()


	//ZK
	c, _, err := zk.Connect([]string{"0.0.0.0:2181"}, 5 * time.Second)
	if err != nil {
		fmt.Printf("Error connecting to zk: %+v", err)
	}

	state := c.State()
	fmt.Printf("state: %+v \n", state)

	createNode := zk.CreateRequest{
		Path: "/aaa",
		Data: []byte{1, 2, 3, 4},
		Acl: zk.WorldACL(zk.PermAll),
		Flags: 0,
		}
	deleteWorkerRqsts := []interface{}{&createNode}
	rep, err := c.Multi(deleteWorkerRqsts...)
	if err != nil {
		fmt.Printf("====  rep! %+v \n", rep)
	}
	//createNode2 := zk.CreateRequest{Path: "/ActiveStreams/s01/node01", Data: []byte{1, 2, 3, 4}, Acl: zk.WorldACL(zk.PermAll)}
	//deleteNode := zk.DeleteRequest{Path: "/ActiveStreams/s01/DONE"}
	//deleteWorkerRqsts := []interface{}{&createNode, &createNode2, &deleteNode}
	//_, err = c.Multi(deleteWorkerRqsts...)
	//if err != nil {
	//	fmt.Printf("====  delete error! %+v \n", err)
	//}

/*	path0 := "/bbb"
	path, err := c.Create(
		path0,
		[]byte{1, 2, 3, 4},
		0,
		zk.WorldACL(zk.PermAll));
	if err != nil {
		fmt.Printf("2 Error : %+v \n", err)
		//panic(err)
	}
	fmt.Printf("== 2 path: %+v \n", path)*/


	children, stat, eventChl, err := c.ChildrenW("/aaa")
	fmt.Printf("==== children: %+v \n", children)
	fmt.Printf("==== stat: %+v \n", stat)
	fmt.Printf("==== err: %+v \n", err)
	for {
		select{
		case event := <-eventChl:
			fmt.Printf("==== watch event: %+v \n", event)

		}
	}

	//a := []string{"le_0000000000", "le_0000000001", "le_0000000002"}
	//b := []string{"le_0000000000", "le_0000000002", "le_0000000001"}
	//fmt.Println(reflect.DeepEqual(a, b))


	/*
	//ZK
	c, _, err := zk.Connect([]string{"192.168.3.130:2181"}, 5 * time.Second) //*10)
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

	*/

	//ZK TEST time out
	/*		start := time.Now()

			/////////// 1
			c, _, err := zk.Connect([]string{"10.0.0.1:2185"}, 10 * time.Second) //*10)
			if err != nil {
				fmt.Printf("1 Error connecting to zk: %+v \n", err)
				//panic(err)
			}
			elapsed1 := time.Since(start)
			fmt.Printf("== 1 elapsed1: %+v \n", elapsed1)


				/////////// 2
				path0 := "/aaaa"
					path, err := c.Create(
						path0,
						[]byte{1, 2, 3, 4},
						0,
						zk.WorldACL(zk.PermAll));
				if err != nil {
					fmt.Printf("2 Error : %+v \n", err)
					//panic(err)
				}
				elapsed2 := time.Since(start)
				fmt.Printf("== 2 elapsed2: %+v \n", elapsed2)
				fmt.Printf("== 2 path: %+v \n", path)



				/////////// 3
				children, _, err := c.Children(path0)
				if err != nil {
					fmt.Printf("3 GET Error :%s, ||  %+v \n", children, err)
				}
				elapsed3 := time.Since(start)
				fmt.Printf("== 3 elapsed2: %+v \n", elapsed3)



				//////////



				///////////
				ch := make(chan int)

				go func(){
					for {
						select {
						case aa := <- ch:
							fmt.Printf("== aa: %+v \n", aa)
						}
					}

				}()

				fmt.Printf("== 4 ch: %+v \n", ch)
				ch <- 1
				close(ch)

				//ch <- 2
				fmt.Printf("== 5 ch: %+v \n", ch)
				//value, ok := <- ch
				//fmt.Printf("== ch: %+v || ok: %v ||value: %v \n", ch, ok, value)
				close(ch)

				select {
				}*/

	//ch := make(chan T)
	//fmt.Println(IsClosed(ch)) // false
	//close(ch)
	//fmt.Println(IsClosed(ch)) // true

	//
	//ch2 := NewMyChannel()
	//
	//ch2.SafeClose()
	////fmt.Println( <- ch2.C)
	//
	//ch2.SafeClose2()
	////fmt.Println( <- ch2.C)



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
type T int

func IsClosed(ch chan T) bool {
	select {
	case some := <-ch:
		fmt.Printf("some: |%+v| \n", some)
		return true
	default:
	}

	return false
}

type MyChannel struct {
	C    chan T
	C2    chan T
	once sync.Once
}

func NewMyChannel() *MyChannel {
	return &MyChannel{
			C: make(chan T),
			C2: make(chan T),
		}
}

func (mc *MyChannel) SafeClose() {
	mc.once.Do(func() {
		fmt.Println("close  C1")
		close(mc.C)
	})
}

func (mc *MyChannel) SafeClose2() {
	mc.once.Do(func() {
		fmt.Println("close  C2")
		close(mc.C2)
	})
}

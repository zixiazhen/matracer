package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"encoding/json"
	"flag"
	"matracer/pkg/api"
	"time"

	rest "gopkg.in/resty.v1"
	//"os"
	//"syscall"
	//"os/signal"
	"fmt"
	"os"
)

const (
	endpointuri      = "/api/v1/namespaces/default/endpoints/"
	ZK_INSTANCE_NUM  = 5
	MA_ENDPOINT_NAME = "zookeeper"
	ZK_PORT          = 2181
)

var (
	apiserver    string
	endpointName string

	//ZK Trace
	traceInterval int
	zkNodeName    string
	zkConns       []*zk.Conn
)

func main() {

	//flags
	flag.StringVar(&apiserver, "apiserver", "http://192.168.3.130:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	flag.StringVar(&endpointName, "endpointname", "zookeeper", "endpoint name, e.g, zookeeper")

	//Trace
	flag.IntVar(&traceInterval, "ma_trace_interval", 5, "watch interval")
	flag.StringVar(&zkNodeName, "zknodename", "/", "zknodename, e.g, /ActiveStream")

	flag.Parse()

	//maEndpointFullPath := apiserver + endpointuri + MA_ENDPOINT_NAME

	//form a list of zk endpoints
	var zkEndpointNamess []string
	for i := 1; i <= ZK_INSTANCE_NUM; i++ {
		zkEndpointNamess = append(zkEndpointNamess, fmt.Sprintf("%s%s%s-0%v", apiserver, endpointuri, MA_ENDPOINT_NAME, i))
	}
	//fmt.Printf("zkEndpointNamess: %v \n", zkEndpointNamess)

	var zkServers []string
	for _, zkEndpointName := range zkEndpointNamess {
		zkServer, _ := getZKFirstEndpoints(zkEndpointName)
		zkServers = append(zkServers, zkServer)
	}

	//This is used for app internal stop signal.
	stop := make(chan error)
	errChal := make(chan error)
	gracefulStop := make(chan os.Signal)

	//Init zookeeper connections
	zkConns := initZKConn(zkServers)
	fmt.Printf("zkConns: %v \n", zkConns)

	//Start trace zk node
	go TraceZK(errChal, gracefulStop)

	//Block
	for {
		select {
		case err := <-stop:
			fmt.Printf(">>>> Stopping the program. Error: %v \n", err.Error())
			return
		case sig := <-gracefulStop:
			fmt.Printf("%v", ">>>> Get a signal to stop the application. \n")
			fmt.Printf(">>>> Stopping the program: %v \n", sig)
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}
	}

	/*
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
			case event := <-ch:
				fmt.Printf("==== [%v] - event: %+v \n", time.Now(), event)
			}
		}*/
}

func initZKConn(zkServers []string) []*zk.Conn {
	var conns []*zk.Conn
	for i := range zkServers {
		c, _, err := zk.Connect([]string{zkServers[i]}, time.Second)
		if err != nil {
			fmt.Printf("Error connecting to zk: %+v", err)
		} else {
			conns = append(conns, c)
		}
	}
	return conns
}

type Event struct {
	stat *zk.Stat
	data []byte
	err  error
}

func TraceZK(errChl chan error, gracefulStop chan os.Signal) {

	// ticker
	ticker := time.NewTicker(time.Duration(traceInterval) * time.Second)
	//quit := make(chan struct{})
	var events chan Event
	for {
		select {
		case e := <-events:
			printEvent(e)
		case <-ticker.C:
			goGetZKStat(errChl, events)
		case <-gracefulStop:
			ticker.Stop()
			fmt.Printf("%v", "Stopping MA Tracer! \n")
			return
		}
	}
}

func printEvent(event Event) {
	fmt.Printf("evnet: %v \n", event)
}

func goGetZKStat(errChl chan error, events chan Event) {
	var event Event
	for i := range zkConns {
		data, stat, err := zkConns[i].Get(zkNodeName)
		if err != nil{
			fmt.Printf("zookeeper connection error,ZK State: %v \n", zkConns[i].State())
			continue
		}
		event.data = data
		event.stat = stat
	}
	events <- event
}

//get the endpoint for 2181 port, which is the first endpoint from the list
func getZKFirstEndpoints(endpointFullPath string) (string, int32) {

	/* Do a simple get from k8s api server */
	resp, err := rest.R().Get(endpointFullPath)
	if err != nil {
		fmt.Printf("Cannot access server: %v \n", err.Error())
		return "", 0
	}

	/* get a list of endpoints */
	//fmt.Printf("%s",resp.Body())
	var eps api.Endpoints
	err = json.Unmarshal(resp.Body(), &eps)
	if err != nil {
		fmt.Print("Unmarshal resp body failed! \n")
		return "", 0
	}

	//var addresses []string
	//need to verify if in all case, there is only one subset in endpoints
	if eps.Subsets == nil || len(eps.Subsets) == 0 || len(eps.Subsets[0].Addresses) == 0 {
		fmt.Printf("No endpoint information found!\n")
		return "", 0
	}

	endpointAddrList := eps.Subsets[0].Addresses
	port := eps.Subsets[0].Ports[0].Port

	fmt.Printf("addr: %v | port: %v \n", endpointAddrList[0].IP, port)

	return endpointAddrList[0].IP, port
}

package main

import (
	"encoding/json"
	"flag"
	"github.com/samuel/go-zookeeper/zk"
	"matracer/pkg/api"
	"time"

	"fmt"
	"github.com/olekukonko/tablewriter"
	rest "gopkg.in/resty.v1"
	"os"
	"os/signal"
	"syscall"
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
	zkConns       map[string]*zk.Conn //zk-IP: Conn
)

func main() {

	//flags
	flag.StringVar(&apiserver, "apiserver", "http://192.168.3.130:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	flag.StringVar(&endpointName, "endpointname", "zookeeper", "endpoint name, e.g, zookeeper")

	//Trace
	flag.IntVar(&traceInterval, "interval", 5, "watch interval")
	flag.StringVar(&zkNodeName, "nodename", "/", "zknodename, e.g, /ActiveStream")

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
	eventChl := make(chan Event)
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGSYS)
	signal.Notify(gracefulStop, syscall.SIGKILL)

	//Init zookeeper connections
	zkConns = initZKConn(zkServers)

	//print result
	go printEvent(eventChl)

	//Start trace zk node
	go TraceZK(errChal, gracefulStop, eventChl)

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
}

type Event struct {
	zkStat []*ZKEvent
}

type ZKEvent struct {
	zkIP string
	stat *zk.Stat
	data []byte
	err  error
}

func TraceZK(errChl chan error, gracefulStop chan os.Signal, eventChl chan Event) {

	// ticker
	ticker := time.NewTicker(time.Duration(traceInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			goGetZKStat(errChl, eventChl)
		case <-gracefulStop:
			ticker.Stop()
			return
		}
	}
}

func printEvent(eventChl chan Event) {
	for {
		select {
		case e := <-eventChl:
			printStats(e.zkStat)
		}
	}
	fmt.Printf(" %v \n", "Stop printing")
}

func printStats(zkStat []*ZKEvent) {
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ZK IP", "Czxid", "Mzxid", "Ctime", "Mtime", "Version", "Cversion", "Aversion", "EphemeralOwner", "DataLength", "NumChildren", "Pzxid"})
	
	//var data [][]string
	for i := range zkStat {
		stat := zkStat[i]
		//fmt.Printf("stat: %v \n" , stat)

		var row []string
		if stat == nil || stat.stat == nil{
			/*
			row = append(row,
				zkStat[i].zkIP,
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
			)
			table.Append(row)
			*/
			continue
		}

		//fmt.Printf("stat.stat: %v \n" , stat.stat)

		row = append(row,
			stat.zkIP,
			fmt.Sprintf("%v", stat.stat.Czxid),
			fmt.Sprintf("%v", stat.stat.Mzxid),
			fmt.Sprintf("%v", stat.stat.Ctime),
			fmt.Sprintf("%v", stat.stat.Mtime),
			fmt.Sprintf("%v", stat.stat.Version),
			fmt.Sprintf("%v", stat.stat.Cversion),
			fmt.Sprintf("%v", stat.stat.Aversion),
			fmt.Sprintf("%v", stat.stat.EphemeralOwner),
			fmt.Sprintf("%v", stat.stat.DataLength),
			fmt.Sprintf("%v", stat.stat.NumChildren),
			fmt.Sprintf("%v", stat.stat.Pzxid),
		)
		//data = append(data, row)
		table.Append(row)
	}
	fmt.Printf("  %v \n" , time.Now())
	table.Render()
}

func goGetZKStat(errChl chan error, eventChl chan Event) {
	var zkStats []*ZKEvent
	for ip, _ := range zkConns {
		var zkStat ZKEvent
		c := zkConns[ip]
		if c == nil {
			fmt.Printf("%v: Cannot connect to ZK.  \n", ip)
			zkStat.zkIP = ip
			zkStat.err = fmt.Errorf("%v", "Cannot connect to ZK!")
			zkStats = append(zkStats, &zkStat)
			continue
		}

		data, stat, err := c.Get(zkNodeName)
		if err != nil {
			fmt.Printf("%v: Get ZK znode failed: %v \n", ip,  err.Error())
			continue
		}
		zkStat.zkIP = ip
		zkStat.data = data
		zkStat.stat = stat
		zkStat.err = nil
		zkStats = append(zkStats, &zkStat)
	}
	event := Event{
		zkStat: zkStats,
	}
	eventChl <- event
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

func initZKConn(zkServers []string) map[string]*zk.Conn {
	//var conns []*zk.Conn
	m := make(map[string]*zk.Conn)
	for i := range zkServers {
		c, _, err := zk.Connect([]string{zkServers[i]}, time.Second)
		if err != nil {
			fmt.Printf("Error connecting to zk: %+v", err)
			if len(zkServers[i]) == 0 {
				continue
			}
			m[zkServers[i]] = nil
		} else {
			m[zkServers[i]] = c
		}
	}
	return m
}

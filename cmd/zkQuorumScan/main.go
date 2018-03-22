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
	"strings"
)

const (
	endpointuri       = "/api/v1/namespaces/default/endpoints/"
	serviceuri        = "/api/v1/namespaces/default/services/"
	ZK_INSTANCE_NUM   = 5
	MA_ENDPOINT_NAME  = "zookeeper"
	ZK_SERVICE_PREFIX = "zookeeper-"
	ZK_PORT           = 2181
)

var (
	apiserver    string
	endpointName string
	zkIPSource   string
	zkServerList  string

	//ZK Trace
	traceInterval int
	showZNodeStat bool
	zkNodeName    string
	zkConns       map[string]*zk.Conn //zk-IP: Conn
	zkServers     []string
)

func main() {

	//flags
	flag.StringVar(&apiserver, "apiserver", "http://192.168.3.130:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	flag.StringVar(&endpointName, "endpointname", "zookeeper", "endpoint name, e.g, zookeeper")
	flag.StringVar(&zkIPSource, "zkipsource", "", "Specify where we get zk ips, e.g, service, endpoint")
	flag.StringVar(&zkServerList, "zkserverlist", "127.0.0.1:2181", "zk ip")

	flag.BoolVar(&showZNodeStat, "traceznode", false, "Enable create and delete stream.")

	//Trace
	flag.IntVar(&traceInterval, "interval", 5, "watch interval")
	flag.StringVar(&zkNodeName, "nodename", "/", "zknodename, e.g, /ActiveStream")

	flag.Parse()

	/*	//form a list of zk endpoints
		var zkEndpointNamess []string
		for i := 1; i <= ZK_INSTANCE_NUM; i++ {
			zkEndpointNamess = append(zkEndpointNamess, fmt.Sprintf("%s%s%s-0%v", apiserver, endpointuri, MA_ENDPOINT_NAME, i))
		}
		//fmt.Printf("zkEndpointNamess: %v \n", zkEndpointNamess)

		for _, zkEndpointName := range zkEndpointNamess {
			zkServer, _ := getZKFirstEndpoints(zkEndpointName)
			zkServers = append(zkServers, zkServer)
		}*/

	//This is used for app internal stop signal.
	stop := make(chan error)
	errChal := make(chan error)
	eventChl := make(chan Event)
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGSYS)
	signal.Notify(gracefulStop, syscall.SIGKILL)

	fmt.Printf("zkIPSource: %s \n", zkIPSource)

	//init zkServers ip
	if zkIPSource == "service" {
		serviceList, err := getZKServices()
		if err != nil {
			fmt.Printf("Get service list failed: err: %s \n", err.Error())
		}
		//GetZKServerIPFromK8SZKService()
		GetZKServerIPFromK8SZKServiceNew(serviceList)
	} else if zkIPSource == "endpoint" {
		GetZKServerIPFromK8SZKEndpoint()
	} else if len(zkServers) > 0 {
		GetZKServerIPFromArgs(zkServerList)
	} else {
		zkServers = []string{"127.0.0.1:2181"}
	}

	//Init zookeeper connections
	zkConns = initZKConn(zkServers)

	//print result
	go printEvent(eventChl)

	//Start tracing zk node
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
	zkStat  []*ZKStatDetails
	svrStat []*zk.ServerStats
}

type ZKStatDetails struct {
	zkIP     string
	stat     *zk.Stat
	data     []byte
	err      error
	children []string
}

func GetZKServerIPFromK8SZKEndpoint() {
	//form a list of zk endpoints
	var zkEndpointNamess []string
	for i := 1; i <= ZK_INSTANCE_NUM; i++ {
		zkEndpointNamess = append(zkEndpointNamess, fmt.Sprintf("%s%s%s-0%v", apiserver, endpointuri, MA_ENDPOINT_NAME, i))
	}
	//fmt.Printf("zkEndpointNamess: %v \n", zkEndpointNamess)

	for _, zkEndpointName := range zkEndpointNamess {
		zkServer, _ := getZKFirstEndpoints(zkEndpointName)
		zkServers = append(zkServers, zkServer)
	}
}

func GetZKServerIPFromK8SZKServiceNew(serviceList api.ServiceList) {
	for _, service := range serviceList.Items {
		if !strings.HasPrefix(service.ObjectMeta.Name, ZK_SERVICE_PREFIX) {
			continue
		}
		zkServer := fmt.Sprintf("%s:%s", service.Spec.ClusterIP, ZK_PORT)
		zkServers = append(zkServers, zkServer)
	}
	fmt.Printf("zkServers: %v \n", zkServers)
}

func GetZKServerIPFromK8SZKService() {
	//form a list of zk service
	var zkServiceNamess []string
	for i := 1; i <= ZK_INSTANCE_NUM; i++ {
		zkServiceNamess = append(zkServiceNamess, fmt.Sprintf("%s%s%s-0%v", apiserver, serviceuri, MA_ENDPOINT_NAME, i))
	}
	fmt.Printf("zkServiceNamess: %v \n", zkServiceNamess)

	for _, zkServicName := range zkServiceNamess {
		zkServer, _ := getZKServicesByFullpath(zkServicName)
		zkServers = append(zkServers, zkServer)
	}
}

// input should look like this "1.2.3.4:2181;2.3.4.5:2181"
func GetZKServerIPFromArgs(zkServerList string) {
	zkServers = strings.Split(zkServerList, ";")
}

func TraceZK(errChl chan error, gracefulStop chan os.Signal, eventChl chan Event) {

	// ticker
	ticker := time.NewTicker(time.Duration(traceInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			//RefreshServerInfoFromEndpoints()
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
			if showZNodeStat == true {
				printNodePathStats(e.zkStat)
			}
			printServerStates(e.svrStat)
		}
	}
	fmt.Printf(" %v \n", "Stop printing")
}

func printServerStates(svrStat []*zk.ServerStats) {
	if len(svrStat) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Sent", "Received", "NodeCount", "MinLatency", "AvgLatency", "MaxLatency", "Connections", "Outstanding", "Epoch", "Counter", "BuildTime", "Mode", "Version", "Error"})

	for i := range svrStat {
		svrStat := svrStat[i]
		var row []string
		if svrStat == nil {
			continue
		}
		row = append(row,
			fmt.Sprintf("%v", svrStat.Sent),
			fmt.Sprintf("%v", svrStat.Received),
			fmt.Sprintf("%v", svrStat.NodeCount),
			fmt.Sprintf("%v", svrStat.MinLatency),
			fmt.Sprintf("%v", svrStat.AvgLatency),
			fmt.Sprintf("%v", svrStat.MaxLatency),
			fmt.Sprintf("%v", svrStat.Connections),
			fmt.Sprintf("%v", svrStat.Outstanding),
			fmt.Sprintf("%v", svrStat.Epoch),
			fmt.Sprintf("%v", svrStat.Counter),
			fmt.Sprintf("%v", svrStat.BuildTime),
			fmt.Sprintf("%v", svrStat.Mode),
			fmt.Sprintf("%v", svrStat.Version),
			fmt.Sprintf("%v", svrStat.Error.Error()),
		)
		table.Append(row)
	}
	table.Render()
}

func printNodePathStats(zkStat []*ZKStatDetails) {

	if len(zkStat) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ZK IP", "Czxid", "Mzxid", "Ctime", "Mtime", "Version", "Cversion", "Aversion", "EphemeralOwner", "DataLength", "NumChildren", "Pzxid", "Children"})

	//Print znode stat
	for i := range zkStat {
		stat := zkStat[i]
		//fmt.Printf("stat: %v \n" , stat)

		var row []string
		if stat == nil || stat.stat == nil {
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
			fmt.Sprintf("%v", stat.children),
		)
		//data = append(data, row)
		table.Append(row)
	}
	fmt.Printf("  %v \n", time.Now())
	table.Render()
}

func goGetZKStat(errChl chan error, eventChl chan Event) {
	//get znode stat
	var zkStats []*ZKStatDetails
	for ip, _ := range zkConns {
		var zkStat ZKStatDetails
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
			fmt.Printf("%v: Get ZK znode failed: %v \n", ip, err.Error())
			continue
		}
		//get children
		children, _, err := c.Children(zkNodeName)
		if err != nil {
			fmt.Printf("%v: Get ZK znode children failed: %v \n", ip, err.Error())
			continue
		}

		zkStat.zkIP = ip
		zkStat.data = data
		zkStat.stat = stat
		zkStat.err = nil
		zkStat.children = children
		zkStats = append(zkStats, &zkStat)
	}

	//get server state
	svrStates, imOK := zk.FLWSrvr(zkServers, 5*time.Second)
	if imOK == false {
		fmt.Printf("ZK server status is innormal! \n")
	}

	//
	event := Event{
		zkStat:  zkStats,
		svrStat: svrStates,
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

//get the endpoint for 2181 port, which is the first service from the list
func getZKServicesByFullpath(serviceFullPath string) (string, int32) {

	/* Do a simple get from k8s api server */
	resp, err := rest.R().Get(serviceFullPath)
	if err != nil {
		fmt.Printf("Cannot access server: %v \n", err.Error())
		return "", 0
	}

	/* get a list of endpoints */
	//fmt.Printf("%s",resp.Body())
	var service api.Service
	err = json.Unmarshal(resp.Body(), &service)
	if err != nil {
		fmt.Print("Unmarshal resp body failed! \n")
		return "", 0
	}

	return service.Spec.ClusterIP, ZK_PORT
}

//get the endpoint for 2181 port, which is the first service from the list
func getZKServices() (api.ServiceList, error) {

	var serviceList api.ServiceList

	serviceUrl := fmt.Sprintf("%s%s", apiserver, serviceuri)
	fmt.Printf("serviceUrl: %+v \n", serviceUrl)

	resp, err := rest.R().Get(serviceUrl)
	if err != nil {
		fmt.Printf("Cannot access server: %v \n", err.Error())
		return serviceList, err
	}

	err = json.Unmarshal(resp.Body(), &serviceList)
	if err != nil {
		fmt.Print("Unmarshal resp body failed! \n")
		return serviceList, err
	}
	fmt.Printf("serviceList: %+v \n", serviceList)

	return serviceList, nil
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

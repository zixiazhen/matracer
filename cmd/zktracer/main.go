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
)

const (
	endpointuri      = "/api/v1/namespaces/default/endpoints/"
	ZK_INSTANCE_NUM  = 5
	MA_ENDPOINT_NAME = "zookeeper"
)

func main() {
	var (
		apiserver    string
		endpointName string
	)

	//flags
	flag.StringVar(&apiserver, "apiserver", "http://127.0.0.1:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	flag.StringVar(&endpointName, "endpointname", "zookeeper", "endpoint name, e.g, zookeeper")

	flag.Parse()

	//form a list of zk endpoints
	var zkEndpointNamess []string
	for i := 1; i <= ZK_INSTANCE_NUM; i++ {
		zkEndpointNamess = append(zkEndpointNamess, fmt.Sprintf("%s-%s", MA_ENDPOINT_NAME, i))
	}

	var zkServers []string
	for _, zkEndpointName := range zkEndpointNamess {
		getZKFirstEndpoints(zkEndpointName)

	}


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
	}
}

//get the endpoint for 2181 port, which is the first endpoint from the list
func getZKFirstEndpoints(endpointFullPath string) (map[string]*api.ObjectReference, int32) {

	//This map cache the endpoint addr for the MA pods.
	endpointsMap := make(map[string]*api.ObjectReference) //IP : Ref

	/* Do a simple get from k8s api server */
	resp, err := rest.R().Get(endpointFullPath)
	if err != nil {
		fmt.Printf("Cannot access server! \n")
		return nil, 0
	}

	/* get a list of endpoints */
	//fmt.Printf("%s",resp.Body())
	var eps api.Endpoints
	err = json.Unmarshal(resp.Body(), &eps)
	if err != nil {
		fmt.Print("Unmarshal resp body failed! \n")
		return nil, 0
	}

	//var addresses []string
	//need to verify if in all case, there is only one subset in endpoints
	if eps.Subsets == nil || len(eps.Subsets) == 0 || len(eps.Subsets[0].Addresses) == 0 {
		fmt.Printf("No endpoint information found! Will quit this program!\n")
		return nil, 0
	}
	endpointAddrList := eps.Subsets[0].Addresses
	port := eps.Subsets[0].Ports[0].Port

	//Get all IPs from endpoints, add to endpointsMap
	for _, endpointAddr := range endpointAddrList {
		endpointsMap[endpointAddr.IP] = endpointAddr.TargetRef
	}

	return endpointsMap, port
}

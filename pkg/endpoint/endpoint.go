package endpoint

import (
	"fmt"
	"encoding/json"
	"matracer/pkg/api"

	rest "gopkg.in/resty.v1"
)


type EndPoint struct{
	endpointFullPath string
	quit chan struct{}
}

func Run(endpointFullPath string, quit chan struct{}) {

	/* Do a simple get from k8s api server */
	resp, err := rest.R().Get(endpointFullPath)
	if err != nil{
		fmt.Printf("Cannot access server!")
	}

	/* get a list of endpoints */
	//fmt.Printf("%s",resp.Body())
	var eps api.Endpoints
	err = json.Unmarshal(resp.Body(), &eps)
	if err != nil {
		fmt.Printf("Unmarshal resp body failed!")
	}

	//var addresses []string
	//need to verify if in all case, there is only one subset in endpoints
	endpointAddrList := eps.Subsets[0].Addresses
	port := eps.Subsets[0].Ports[0].Port
	fmt.Print(port)

	for _, endpointAddr := range endpointAddrList {
		fmt.Printf(endpointAddr.IP)
		//addresses = append(addresses, addr)
	}


}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"matracer/pkg/api"
	"time"

	rest "gopkg.in/resty.v1"
)

const (
	endpointuri = "/api/v1/namespaces/default/endpoints/"
)

func main() {
	var (
		apiserver    string
		endpointName string
		frequency int
	)

	/* Handling flags */
	flag.StringVar(&apiserver, "apiserver", "http://127.0.0.1:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	flag.StringVar(&endpointName, "endpointname", "manifest-agent", "endpoint name, e.g, manifest-agent")
	flag.IntVar(&frequency, "frequency", 5, "watch frequency")
	flag.Parse()

	/* endpoint full url */
	endpointFullPath := apiserver + endpointuri + endpointName
	fmt.Printf("endpointFullPath: %s \n",endpointFullPath)

	/* ticker */
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	//quit := make(chan struct{})
	stop := make(chan error)

	for {
		select {
		case <-ticker.C:
			run(endpointFullPath, stop)
		case err := <-stop:
			ticker.Stop()
			fmt.Printf("Error: %v", err.Error())
			return
		}
	}
}

func run(endpointFullPath string, stop chan error) {

	//This map cache the endpoint addr for the MA pods.
	endpointsMap := make(map[string]*api.ObjectReference) //IP : Ref

	//This map is use as a set. to determine if there are duplicates
	streamMap := make(map[string]string) //StreamID : IP

	//This map is used for show the result
	// MA that own a stream:	PodName : "StreamID"
	// MA that is idle:			PodName : "idle"
	result := make(map[string]string)

	/* Do a simple get from k8s api server */
	resp, err := rest.R().Get(endpointFullPath)
	if err != nil {
		fmt.Printf("Cannot access server! \n")
		return
	}

	/* get a list of endpoints */
	//fmt.Printf("%s",resp.Body())
	var eps api.Endpoints
	err = json.Unmarshal(resp.Body(), &eps)
	if err != nil {
		fmt.Print("Unmarshal resp body failed! \n")
		return
	}

	//var addresses []string
	//need to verify if in all case, there is only one subset in endpoints
	if eps.Subsets == nil || len(eps.Subsets) == 0 || len(eps.Subsets[0].Addresses) == 0 {
		fmt.Printf("No endpoint information found! Will quit this program!\n")
		return
	}
	endpointAddrList := eps.Subsets[0].Addresses
	port := eps.Subsets[0].Ports[0].Port

	//Get all IPs from endpoints, add to endpointsMap
	for _, endpointAddr := range endpointAddrList {
		endpointsMap[endpointAddr.IP] = endpointAddr.TargetRef
	}

	//range endpointsMap
	for epIPAddr, epObjRef := range endpointsMap {
		if epObjRef == nil {
			fmt.Print("Object Reference is nil, skip this pod! \n")
			continue
		}

		//1. Create the MA status rest call url
		maStatusRestCall := fmt.Sprintf("http://%s:%v/status", epIPAddr, port)
		fmt.Printf("maStatusRestCall: %s \n",maStatusRestCall)

		//2. Get steam ID from MA Status
		maStatusRaw, err := rest.R().Get(maStatusRestCall)
		if err != nil {
			fmt.Printf("Cannot access Manifest Agent: %s \n", epIPAddr)
			result[epObjRef.Name] = "Down"
			continue
		}

		var maStatus api.MAStatus
		err = json.Unmarshal(maStatusRaw.Body(), &maStatus)
		if err != nil {
			fmt.Print("Unmarshal ma status resp body failed!")
			result[epObjRef.Name] = "Status Unknown"
			continue
		}

		//3. Check if there are multiple MA hold the same Stream ID
		if len(maStatus.StreamID) != 0 {
			//this pod own a stream.
			//Check if there is another endpoint already own the stream.
			if anotherEndpoint, found := streamMap[maStatus.StreamID]; found {
				fmt.Printf(" =======================================================  \n")
				fmt.Printf(" ============ Multi-MA own the same steam!! ============  \n")
				fmt.Printf(" | Strean ID: 	%v  \n", maStatus.StreamID)
				fmt.Printf(" | MA-1: 		%v  \n", epObjRef.Name)
				fmt.Printf(" | MA-2: 		%v  \n", endpointsMap[anotherEndpoint].Name)
				fmt.Printf(" =======================================================  \n")

				//add a record result
				result[epObjRef.Name] = maStatus.StreamID + "****"
				result[endpointsMap[anotherEndpoint].Name] = maStatus.StreamID + "****"
			} else {
				streamMap[maStatus.StreamID] = epIPAddr
				//add a record result
				result[epObjRef.Name] = maStatus.StreamID
			}
		} else {
			//add a record result
			result[epObjRef.Name] = "Idle"
		}
	}

	//Print result
	printResult(result)
}

func printResult(result map[string]string) {

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

	//Print the result
	fmt.Printf(" ================ %v ===============  \n", time.Now())
	fmt.Print(string(b))

}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"matracer/pkg/api"
	stream "matracer/pkg/streamcreater"
	"time"

	rest "gopkg.in/resty.v1"
	"os"
)

var (
	apiserver    string
	endpointName string
	action string

	//MA Trace
	enableTrace bool
	maTraceFrequency int

	//Create and delete stream
	enableStreamCreation bool
	streamCreateDeleteFrequency int
	numOfStream int

	//Scale RC
	enableRCScaling bool
)

const (
	endpointuri = "/api/v1/namespaces/default/endpoints/"
	MA_ENDPOINT_NAME = "manifest-agent"
	NSA_ENDPOINT_NAME = "nsa"
)

func main() {

	//flags
	flag.StringVar(&apiserver, "apiserver", "http://127.0.0.1:8080", "url for k8s api server, e.g., http://127.0.0.1:8080")
	//flag.StringVar(&endpointName, "endpointname", "manifest-agent", "endpoint name, e.g, manifest-agent")
	flag.StringVar(&action, "type", "traceonly", "support create/deleteall")

	//Trace
	flag.BoolVar(&enableTrace, "enabletrace", true, "This will keep tracing MA status.")
	flag.IntVar(&maTraceFrequency, "matracefrequency", 5, "watch frequency")

	//Create and delete stream
	flag.BoolVar(&enableStreamCreation, "enablestream", false, "Enable create and delete stream.")
	flag.IntVar(&numOfStream, "numofstream", 2, "The number of stream that will be created.")
	flag.IntVar(&streamCreateDeleteFrequency, "streamperiod", 10, "The period that we create and delete stream.")

	//Scale RC
	flag.BoolVar(&enableRCScaling, "enablescale", false, "This will scale rc in and out.")

	flag.Parse()

	stop := make(chan error)
	errChal := make(chan error)

	if enableTrace == true {
		go traceMA(errChal, stop)
	}

	if enableStreamCreation == true {
		go createAndDeleteStream(errChal)
	}

	if enableRCScaling == true {
		go scaleRC(errChal)
	}

	//Block
	for {
		select{
		case err := <-stop:
			fmt.Printf("%v", ">>>> This should not be reached. \n")
			fmt.Printf("Stopping the program. Error: %v", err.Error())
			return
		}
	}
}


func traceMA(errChl chan error, stop chan error){

	//end point full path
	maEndpointFullPath := apiserver + endpointuri + MA_ENDPOINT_NAME
	// ticker
	ticker := time.NewTicker(time.Duration(maTraceFrequency) * time.Second)
	//quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			fmt.Printf("%v", ">>>> run goTraceMA() \n")
			goTraceMA(maEndpointFullPath, stop, errChl)
		case err := <-stop:
			ticker.Stop()
			fmt.Printf("Error: %v", err.Error())
			return
		}
	}
}

func createAndDeleteStream(errChl chan error) {

	//end point full path
	nsaEndpointFullPath := apiserver + endpointuri + NSA_ENDPOINT_NAME
	//ticker
	ticker := time.NewTicker(time.Duration(streamCreateDeleteFrequency) * time.Second)

	for {
		select {
		case <-ticker.C:
			goCreateAndDeleteStream(nsaEndpointFullPath)
		case err := <-errChl:
			ticker.Stop()
			fmt.Printf(">>>> Got an error from MA: %v \n", err.Error())

			//collectLogs()
			break
		}
	}
}

func goCreateAndDeleteStream(nsaEndpointFullPath string){

	endpointsMap, port := getEndpoints(nsaEndpointFullPath)

	//there should be only one NSA
	for epIPAddr, epObjRef := range endpointsMap {
		if epObjRef == nil {
			fmt.Print("Object Reference is nil, skip this pod! \n")
			continue
		}

		//1. Create the rest call url
		nsaRestCall := fmt.Sprintf("http://%s:%v", epIPAddr, port)
		//fmt.Printf("maStatusRestCall: %s \n",maStatusRestCall)

		//2. Create streams
		stream.Create(nsaRestCall, numOfStream)

		//3. Sleep for a while
		time.Sleep(time.Duration(streamCreateDeleteFrequency) * time.Second)

		//4. Delete all streams
		stream.DeleteAll(nsaRestCall)
	}



}

func collectLogs() {
	fmt.Printf("%v",">>>> Start to collect log. \n")
	defer fmt.Printf("%v",">>>> Finish collecting log. \n")
	return
	}

func goTraceMA(endpointFullPath string, stop chan error, errChl chan error) {

	endpointsMap, port := getEndpoints(endpointFullPath)

	//This map is use as a set. to determine if there are duplicates
	streamMap := make(map[string]string) //StreamID : IP

	//This map is used for show the result
	// MA that own a stream:	PodName : "StreamID"
	// MA that is idle:			PodName : "idle"
	result := make(map[string]string)

	//range endpointsMap
	for epIPAddr, epObjRef := range endpointsMap {
		if epObjRef == nil {
			fmt.Print("Object Reference is nil, skip this pod! \n")
			continue
		}

		//1. Create the MA status rest call url
		maStatusRestCall := fmt.Sprintf("http://%s:%v/status", epIPAddr, port)
		//fmt.Printf("maStatusRestCall: %s \n",maStatusRestCall)

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
			fmt.Print("Unmarshal ma status resp body failed! \n")
			result[epObjRef.Name] = "Status Unknown"
			continue
		}

		//3. Check if there are multiple MA hold the same Stream ID
		if len(maStatus.StreamID) != 0 {
			//this pod own a stream.
			//Check if there is another endpoint already own the stream.
			if anotherEndpoint, found := streamMap[maStatus.StreamID]; found {
				fmt.Printf("\n")
				fmt.Printf(" |------------------------------------------------------|  \n")
				fmt.Printf(" |----------- Multi-MA hold the same stream! -----------|  \n")
				fmt.Printf(" |------------------------------------------------------|  \n")
				fmt.Printf(" | Strean ID:	%v		\n", maStatus.StreamID)
				fmt.Printf(" | MA-1:	%v		\n", epObjRef.Name)
				fmt.Printf(" | MA-2:	%v		\n", endpointsMap[anotherEndpoint].Name)
				fmt.Printf(" |------------------------------------------------------|  \n")
				fmt.Printf("\n")

				//add a record result
				result[epObjRef.Name] = maStatus.StreamID + "	-*-"
				result[endpointsMap[anotherEndpoint].Name] = maStatus.StreamID + "	-*-"

				//Notify stream creater to stop
				errChl <- fmt.Errorf("%v","Multiple MA hold the same stream! \n")

				//create a error file in /tmp folder
				createFile()

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


func getEndpoints(endpointFullPath string) (map[string]*api.ObjectReference, int32) {

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


var path = "/tmp/maerror.log"
func createFile() {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			fmt.Println("err : %v \n", err)
			return
			}
		defer file.Close()
		fmt.Println("==> done creating file \n", path)
	}
}

func printResult(result map[string]string) {

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	//Print the result
	fmt.Printf("\n |------- %v ------|  \n", time.Now())
	//fmt.Printf(" ================ %v ===============  \n", time.Now())
	fmt.Print(string(b))
	fmt.Printf("\n |------------------------------------------------------|  \n")

}



func scaleRC(errChl chan error) {
	fmt.Printf("Coming Soon! \n")
}


/*
func runOld(endpointFullPath string, stop chan error) {

	//This map cache the endpoint addr for the MA pods.
	endpointsMap := make(map[string]*api.ObjectReference) //IP : Ref

	//This map is use as a set. to determine if there are duplicates
	streamMap := make(map[string]string) //StreamID : IP

	//This map is used for show the result
	// MA that own a stream:	PodName : "StreamID"
	// MA that is idle:			PodName : "idle"
	result := make(map[string]string)

	// Do a simple get from k8s api server
	resp, err := rest.R().Get(endpointFullPath)
	if err != nil {
		fmt.Printf("Cannot access server! \n")
		return
	}

	// get a list of endpoints
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
		//fmt.Printf("maStatusRestCall: %s \n",maStatusRestCall)

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

				//create a error file in /tmp folder
				createFile()
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

*/


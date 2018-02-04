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
	"syscall"
	"os/signal"
)

var (
	apiserver    string
	endpointName string
	action string

	//MA Trace
	maTraceEnable   bool
	maTraceInterval int

	//Create and delete stream
	streamAddDelEnable   bool
	streamAddDelInterval int
	streamAddNum         int

	//Scale RC
	scaleEnable bool
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
	//flag.StringVar(&action, "type", "traceonly", "support create/deleteall")

	//Trace
	flag.BoolVar(&maTraceEnable, "ma_trace_enable", false, "This will keep tracing MA status.")
	flag.IntVar(&maTraceInterval, "ma_trace_interval", 5, "watch interval")

	//Create and delete stream
	flag.BoolVar(&streamAddDelEnable, "stream_add_del_enable", false, "Enable create and delete stream.")
	flag.IntVar(&streamAddNum, "stream_add_num", 2, "The number of stream that will be created.")
	flag.IntVar(&streamAddDelInterval, "stream_add_interval", 60, "The interval that we create and delete stream.")

	//Scale RC
	flag.BoolVar(&scaleEnable, "scale_enable", false, "This will scale rc in and out.")

	flag.Parse()

	//This is a stop signal that from os, when we kill the application from os, the app should do clean up and then stop.
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGSYS)
	signal.Notify(gracefulStop, syscall.SIGKILL)

	//This is used for app internal stop signal.
	stop := make(chan error)
	//This channel is used by MA tracer to notify other components to stop.
	//For example, when we observe that two MA hold the same stream, we should notify streamCreater to stop creating stream.
	errChal := make(chan error)

	if maTraceEnable == true {
		go TraceMA(errChal, gracefulStop)
	}

	if streamAddDelEnable == true {
		go CreateAndDeleteStream(errChal, gracefulStop)
	}

	if scaleEnable == true {
		go scaleRC(errChal, gracefulStop)
	}


	//Block
	for {
		select{
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

func TraceMA(errChl chan error, gracefulStop chan os.Signal){

	//end point full path
	maEndpointFullPath := apiserver + endpointuri + MA_ENDPOINT_NAME
	// ticker
	ticker := time.NewTicker(time.Duration(maTraceInterval) * time.Second)
	//quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			goTraceMA(maEndpointFullPath, errChl)
		case <-gracefulStop:
			ticker.Stop()
			fmt.Printf("%v", "Stopping MA Tracer! \n")
			return
		}
	}
}

func CreateAndDeleteStream(errChl chan error, gracefulStop chan os.Signal) {

	//end point full path
	nsaEndpointFullPath := apiserver + endpointuri + NSA_ENDPOINT_NAME
	//ticker
	ticker := time.NewTicker(time.Duration(streamAddDelInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			goCreateAndDeleteStream(nsaEndpointFullPath)
		case err := <-errChl:
			ticker.Stop()
			fmt.Printf(">>>> Got an error from MA: %v \n", err.Error())

			//collectLogs()
			break
		case <- gracefulStop:
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
		streams, err := stream.Create(nsaRestCall, streamAddNum)
		if err != nil {
			fmt.Printf("Create stream failed: %v \n", err.Error())
		}
		//Sleep for a while, wait for the creation done.
		time.Sleep(30 * time.Second)

		//3. Delete all streams
		//stream.DeleteAll(nsaRestCall)
		err = stream.Delete(nsaRestCall, streams)
		if err != nil {
			fmt.Printf("Delete stream failed: %v \n", err.Error())
		}
		//Sleep for a while, wait for the deletion done.
		time.Sleep(20 * time.Second)
	}
}

func collectLogs() {
	fmt.Printf("%v",">>>> Start to collect log. \n")
	defer fmt.Printf("%v",">>>> Finish collecting log. \n")
	return
}

func goTraceMA(endpointFullPath string, errChl chan error) {

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
				result[epObjRef.Name] = maStatus.StreamID + "----"
				result[endpointsMap[anotherEndpoint].Name] = maStatus.StreamID + "----"

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
	fmt.Print(string(b))
	fmt.Printf("\n |------------------------------------------------------|  \n")

}

func scaleRC(errChl chan error, stop chan os.Signal) {
	fmt.Printf("Coming Soon! \n")
}

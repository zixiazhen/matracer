package main

import (
	"flag"
	"fmt"
	"time"
	"matracer/pkg/api"
	rest "gopkg.in/resty.v1"
	"math/rand"
	"encoding/json"
)

const (
	endpointuri = "/api/v1/namespaces/default/endpoints/"
	CreateStreans ACTION = "create"
	GetALlStreams ACTION = "getall"
	DeleteALlStreams ACTION = "deleteall"
	high = 999999999999999999
	low =  100000000000000000
)

type ACTION string

func init(){
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		nsa    string
		num 	int
		action string
	)

	/* Handling flags */
	flag.StringVar(&nsa, "nsa", "http://10.10.224.29:9013", "url for nsa server, e.g., http://10.10.224.29:9013")
	flag.IntVar(&num, "num", 5, "watch frequency")
	flag.StringVar(&action, "action", "create", "create streams")
	flag.Parse()

	//For now, only support create.
	Create(nsa, num)
}

func Create(nsa string, num int){

	fmt.Printf("Start Creating \n")
	defer fmt.Printf("Finish Creating \n")

	/* endpoint full url */
	fullPath := nsa + "/addstreams"

	streams := []api.StreamCfg{}
	for i := 0; i < num; i++ {
		streamConf := generateStream()
		streams = append(streams, streamConf)
	}
	fmt.Printf("streamConf: %v \n", streams)

	bytes, _ := json.Marshal(streams)
	fmt.Printf("stream json: \n %v \n", string(bytes))

	resp, err :=  rest.R().
		SetHeader("Content-Type", "application/json").
		SetBody(streams).
		Post(fullPath)
	if err != nil {
		fmt.Printf("Create Stream failed! %v \n", err.Error())
		return
	}

	fmt.Printf("resp: %v \n", resp)
}

func Delete(){

}

func Get() {

}

func generateStream() api.StreamCfg {

	num := low + rand.Intn(high - low)
	transports := [1]api.TransportCfg{{
		URL: "http://ccr.linear-nat-dash.xcr.comcast.net/dash/USA_SD_NAT_4183_0_7503892744946620183/USA_SD_NAT_4183_0_7503892744946620163_DASH.mpd",
		AvgBitrate:1.875,
		MaxBitrate:1.875,
	}}

	ac := api.ArchiveCfg{
		ArchiveTime: 72,
		ReArchiveTime: 24,
		ArchivalStartTime: "05:00PM",
		ArchivalDuration: 10,
		ArchivalPause: 5,
	}

	stream := api.StreamCfg{
		StreamName: fmt.Sprintf("USA_SD_NAT_4184_%v", num),
		StreamID: fmt.Sprintf("%v", num),
		Expiry: "2018-02-04T20:00:00Z",
		Transports: transports,
		ArchiveConfig:ac,
	}

	return stream
}

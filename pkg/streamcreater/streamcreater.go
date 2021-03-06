package streamcreater


import (
	"fmt"
	"matracer/pkg/api"
	rest "gopkg.in/resty.v1"
	"math/rand"
	//"encoding/json"
	"encoding/json"
)

const (
	high = 999999999999999999
	low =  100000000000000000
)

func Create(nsa string, num int) ([]api.StreamCfg, error) {

	/* endpoint full url */
	fullPath := nsa + "/addstreams"

	streams := []api.StreamCfg{}
	for i := 0; i < num; i++ {
		streamConf := generateStream()
		streams = append(streams, streamConf)
	}

	resp, err :=  rest.R().
		SetHeader("Content-Type", "application/json").
		SetBody(streams).
		Post(fullPath)
	if err != nil {
		fmt.Printf("Create Stream failed! %v \n", err.Error())
		return nil, err
	}

	fmt.Printf(">>>> Created %v streams. Status code: %v \n", num, resp.StatusCode())
	return streams, nil

	}

func Delete(nsa string, streams []api.StreamCfg) error {

	fullPath := nsa + "/configuration"
	for i, _ := range streams{
		url := fmt.Sprintf("%v/%v", fullPath, streams[i].StreamID)
		delResp, err :=  rest.R().
			SetHeader("Content-Type", "application/json").
			Delete(url)
		if err != nil {
			fmt.Printf(">>>> Delete Stream failed! %v \n", err.Error())
			return err
		}
		fmt.Printf(">>>> Deleted stream %v. Status code: %v \n",streams[i].StreamID, delResp.StatusCode())
	}
	return nil
}

func DeleteAll(nsa string) error{

	fullPath := nsa + "/configuration"
	resp, err :=  rest.R().
		SetHeader("Content-Type", "application/json").
		Get(fullPath)
	if err != nil {
		fmt.Printf("Get Stream failed! %v \n", err.Error())
		return err
	}

	streams := make([]api.StreamCfg,0)
	json.Unmarshal(resp.Body(), &streams)

	for i, _ := range streams{
		url := fmt.Sprintf("%v/%v", fullPath, streams[i].StreamID)
		delResp, err :=  rest.R().
			SetHeader("Content-Type", "application/json").
			Delete(url)
		if err != nil {
			fmt.Printf(">>>> Delete Stream failed! %v \n", err.Error())
			return err
		}
		fmt.Printf(">>>> del all streams. Status code: %v \n", delResp.StatusCode())
	}
	return nil
}

func Get() {

}

func generateStream() api.StreamCfg {

	num := low + rand.Intn(high - low)
	transports := [1]api.TransportCfg{{
		URL: "http://172.22.111.238/starzkids_pillar02/manifest.mpd",
		//URL: "http://ccr.linear-nat-dash.xcr.comcast.net/dash/USA_SD_NAT_4183_0_7503892744946620183/USA_SD_NAT_4183_0_7503892744946620163_DASH.mpd",
		AvgBitrate:1.875,
		MaxBitrate:1.875,
	}}

	ac := api.ArchiveCfg{
		ArchiveTime: 30,
		ReArchiveTime: 30,
		ArchivalStartTime: "05:00PM",
		ArchivalDuration: 86400,
		ArchivalPause: 5000,
	}

	stream := api.StreamCfg{
		StreamName: fmt.Sprintf("USA_SD_NAT_4184_%v", num),
		StreamID: fmt.Sprintf("%v", num),
		Expiry: "2018-06-04T20:00:00Z",
		Transports: transports,
		ArchiveConfig:ac,
	}

	return stream
}

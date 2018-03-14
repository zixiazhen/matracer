package main

import (
	//"context"

	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/api/types/mount"
	//"github.com/docker/docker/api/types/volume"
	//"fmt"
	//"matracer/pkg/zkClusterBuilder"
	//"matracer/pkg/zkClusterBuilder"
	//"matracer/pkg/zkClusterBuilder"
	"os"
	"fmt"
)

const (
	ZK_IMAGE = "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test"
)

func main() {
	//zkClusterBuilder.RunRioZookeeperCluster()

	idFile := "/tmp/1/2/3/conf"
	//fileName := "/conf"

	var err error
	_, err = os.Stat(idFile)
	if os.IsNotExist(err) {
		err = os.MkdirAll(idFile, os.ModePerm)
	}
	if err != nil {
		fmt.Printf("Error MkdirAll : %s\n", err)
	}


}

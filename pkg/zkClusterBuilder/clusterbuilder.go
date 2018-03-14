package zkClusterBuilder

import (
	"fmt"
	"runtime"

	"github.com/samalba/dockerclient"
	//"context"

	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/api/types/mount"
	//"github.com/docker/docker/api/types/volume"

	//"github.com/docker/go-connections/nat"
)

const (
	HostAddress = "0.0.0.0"
	ZK_IMAGE    = "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test"
)

type ZKCluster struct {
	members []string
	conf    ZKClusterConf
}

type ZKMember struct {
	address      string
	myid         string
	clientPort   string
	followerPort string
	electionPort string
}

type ZKClusterConf struct {
	numOfMember        string
	containerImageName string
	env                []string
}

func NewZKCluster(conf ZKClusterConf) *ZKCluster {

	return nil
}

func (m *ZKCluster) run() {

}

func (m *ZKCluster) startCluster() {

}

func (m *ZKCluster) stopCluster() {

}

func (m *ZKCluster) initClusterConfig() {
}

func (m *ZKCluster) getMembers() []string {
	return m.members
}

func RunRioZookeeperCluster() {

	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return
	}

	// Stop any previous container
	_ = docker.RemoveContainer("rio-zk", true, true)
	_ = docker.RemoveContainer("rio-zk1", true, true)
	_ = docker.RemoveContainer("rio-zk2", true, true)
	_ = docker.RemoveContainer("rio-zk3", true, true)

	//ZK 1
	conf1 := &dockerclient.ContainerConfig{
		Image: "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test",
		Env: []string{
			"GLOBAL_OUTSTANDING_LIMIT=1000",
			"ZK_MAX_SESSION_TIMEOUT=60000",
			"ZK_INIT_LIMIT=100",
			"ZK_JMX_PORT=9041",
			"ZK_CNX_TIMEOUT=10",
			"ZK_TIMEOUT=5s",
			"MAX_CLIENT_CONNS=5000",
			"LEADER_SERVES=yes",
			"SKIP_ACL=yes",
			"ACTIVE_PATH=gemini-bigs-active",
			"SITE_ID=Cisco",
			"STORAGE_TIMEOUT=10",

			"MYID=1",
			"CLIENT_PORT=2181",
			"ZOOKEEPER_01_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_02_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_03_SERVICE_HOST=" + HostAddress,
			"VERSION=03012018",
		},
	}

	//ZK 2
	conf2 := &dockerclient.ContainerConfig{
		Image: "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test",
		Env: []string{
			"GLOBAL_OUTSTANDING_LIMIT=1000",
			"ZK_MAX_SESSION_TIMEOUT=60000",
			"ZK_INIT_LIMIT=100",
			"ZK_JMX_PORT=9041",
			"ZK_CNX_TIMEOUT=10",
			"ZK_TIMEOUT=5s",
			"MAX_CLIENT_CONNS=5000",
			"LEADER_SERVES=yes",
			"SKIP_ACL=yes",
			"ACTIVE_PATH=gemini-bigs-active",
			"SITE_ID=Cisco",
			"STORAGE_TIMEOUT=10",

			"MYID=2",
			"CLIENT_PORT=2182",
			"ZOOKEEPER_01_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_02_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_03_SERVICE_HOST=" + HostAddress,
			"VERSION=03012018",
		},
	}

	//ZK 3
	conf3 := &dockerclient.ContainerConfig{
		Image: "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test",
		Env: []string{
			"GLOBAL_OUTSTANDING_LIMIT=1000",
			"ZK_MAX_SESSION_TIMEOUT=60000",
			"ZK_INIT_LIMIT=100",
			"ZK_JMX_PORT=9041",
			"ZK_CNX_TIMEOUT=10",
			"ZK_TIMEOUT=5s",
			"MAX_CLIENT_CONNS=5000",
			"LEADER_SERVES=yes",
			"SKIP_ACL=yes",
			"ACTIVE_PATH=gemini-bigs-active",
			"SITE_ID=Cisco",
			"STORAGE_TIMEOUT=10",

			"MYID=3",
			"CLIENT_PORT=2183",
			"ZOOKEEPER_01_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_02_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_03_SERVICE_HOST=" + HostAddress,
			"VERSION=03012018",
		},
	}

	zkID1, err := docker.CreateContainer(conf1, "rio-zk1", nil)
	if err != nil {
		fmt.Printf("==== CreateContainer err: %v", err)
		return
	}
	fmt.Printf("==== zkID1: %v", zkID1)

	zkID2, err := docker.CreateContainer(conf2, "rio-zk2", nil)
	if err != nil {
		fmt.Printf("==== CreateContainer err: %v", err)
		return
	}
	fmt.Printf("==== zkID2: %v", zkID2)

	zkID3, err := docker.CreateContainer(conf3, "rio-zk3", nil)
	if err != nil {
		fmt.Printf("==== CreateContainer err: %v", err)
		return
	}
	fmt.Printf("==== zkID3: %v", zkID3)


	// Start the container
	hostConfig1 := &dockerclient.HostConfig{
		PortBindings: map[string][]dockerclient.PortBinding{
			"2181/tcp": {{"", "2181"}},
			"2888/tcp": {{"", "2888"}},
			"3888/tcp": {{"", "3888"}},
		},
		Binds: []string{
			"/tmp:/tmp",
		},
	}
	// Start the container
	hostConfig2 := &dockerclient.HostConfig{
		PortBindings: map[string][]dockerclient.PortBinding{
			"2182/tcp": {{"", "2182"}},
			"2889/tcp": {{"", "2889"}},
			"3889/tcp": {{"", "3889"}},
		},
		Binds: []string{
			"/tmp:/tmp",
		},
	}
	// Start the container
	hostConfig3 := &dockerclient.HostConfig{
		PortBindings: map[string][]dockerclient.PortBinding{
			"2183/tcp": {{"", "2183"}},
			"2890/tcp": {{"", "2890"}},
			"3890/tcp": {{"", "3890"}},
		},
		Binds: []string{
			"/tmp:/tmp",
		},
	}

	err = docker.StartContainer(zkID1, hostConfig1)
	if err != nil {
		fmt.Printf("==== StartContainer err: %v", err)
		return
	}

	err = docker.StartContainer(zkID2, hostConfig2)
	if err != nil {
		fmt.Printf("==== StartContainer err: %v", err)
		return
	}

	err = docker.StartContainer(zkID3, hostConfig3)
	if err != nil {
		fmt.Printf("==== StartContainer err: %v", err)
		return
	}


	// Get contaienr IP


}

func RunRioZookeeperCluster_bak() {

	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return
	}

	// Stop any previous container
	_ = docker.RemoveContainer("rio-zk", true, true)

	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image: "dockerhub.cisco.com/vmr-docker/rio-zk/zk:test",
		Env: []string{
			"GLOBAL_OUTSTANDING_LIMIT=1000",
			"ZK_MAX_SESSION_TIMEOUT=60000",
			"ZK_INIT_LIMIT=100",
			"ZK_JMX_PORT=9041",
			"ZK_CNX_TIMEOUT=10",
			"ZK_TIMEOUT=5s",
			"MAX_CLIENT_CONNS=5000",
			"LEADER_SERVES=yes",
			"SKIP_ACL=yes",
			"ACTIVE_PATH=gemini-bigs-active",
			"SITE_ID=Cisco",
			"STORAGE_TIMEOUT=10",
			"MYID=1",
			"CLIENT_PORT=2188",
			"ZOOKEEPER_01_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_02_SERVICE_HOST=" + HostAddress,
			"ZOOKEEPER_03_SERVICE_HOST=" + HostAddress,
			"VERSION=03012018",
		},
	}

	zkID, err := docker.CreateContainer(containerConfig, "rio-zk", nil)
	if err != nil {
		fmt.Printf("==== CreateContainer err: %v", err)
		return
	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		PortBindings: map[string][]dockerclient.PortBinding{
			"2181/tcp": {{"", "2181"}},
			"2888/tcp": {{"", "2888"}},
			"3888/tcp": {{"", "3888"}},
		},
	}
	err = docker.StartContainer(zkID, hostConfig)
	if err != nil {
		fmt.Printf("==== StartContainer err: %v", err)
		return
	}

}

func ContainerIP(ID string, dc *dockerclient.DockerClient) (string, error) {
	if runtime.GOOS == "darwin" {
		// Docker for Mac doesn't support per-container IP addressing, so just use localhost instead.
		return "127.0.0.1", nil
	}

	res, err := dc.InspectContainer(ID)
	if err != nil {
		return "", fmt.Errorf("could not inspect docker container %s", err)
	}
	return res.NetworkSettings.IPAddress, nil
}

//StopZookeeper - This API stops the running ZK container
//
// Return values:
// err error - On success returns nil, on failure a non-nil error object shall be sent
func StopZookeeper() error {
	return StopContainer("rio-zk-test")
}

//StopContainer - Stop a running container identified by the id.
//
// Return values:
// err error - On success returns nil, on failure a non-nil error object shall be sent
func StopContainer(id string) error {
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return err
	}
	if err := docker.StopContainer(id, 2); err != nil {
		return err
	}
	return nil
}


package discovery

import "github.com/jiuzhou-zhao/go-fundamental/clienttoolset"

type ServiceInfo struct {
	Name          string
	GRpcAddresses *clienttoolset.GRpcClientConfig
}

type Observer func(services map[string][]*ServiceInfo)

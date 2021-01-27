package discovery

import (
	"fmt"
	"strings"
)

const (
	TypeGRpc = "grpc"
	TypeHttp = "http"
)

type ServiceInfo struct {
	Host        string
	Port        int
	ServiceName string // type:name:index
	Meta        map[string]string
}

type Observer func(services map[string]*ServiceInfo)

func BuildServerName(t string, name string, index string) string {
	t = strings.ToLower(t)
	return fmt.Sprintf("%v:%v:%v", t, name, index)
}

func ParseServerName(n string) (t string, name string, index string, err error) {
	vs := strings.Split(n, ":")
	if len(vs) < 2 || len(vs) > 3 {
		err = fmt.Errorf("invalid server name: %v", n)
		return
	}
	t = vs[0]
	t = strings.ToLower(t)
	name = vs[1]
	if len(vs) == 3 {
		index = vs[2]
	}
	return
}

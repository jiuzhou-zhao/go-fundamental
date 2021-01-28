package discovery

import (
	"fmt"
	"strings"
)

func BuildServerName(t string, name string, index string) string {
	t = strings.ToLower(t)
	serverName := fmt.Sprintf("%v:%v", t, name)
	if index != "" {
		serverName = serverName + ":" + index
	}
	return serverName
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

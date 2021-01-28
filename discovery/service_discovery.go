package discovery

const (
	TypeGRpc = "grpc"
	TypeHttp = "http"
)

const (
	MetaGRPCClass = "grpc_class"
)

type ServiceInfo struct {
	Host        string
	Port        int
	ServiceName string // type:name:index
	Meta        map[string]string
}

type Observer func(services []*ServiceInfo)

//
//
//
type Setter interface {
	Start([]*ServiceInfo) error
	Stop()
	Wait()
}

type Getter interface {
	Start(ob Observer, opt ...Option) error
	Stop()
	Wait()
}

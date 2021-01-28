package servicetoolset

import (
	"github.com/jiuzhou-zhao/go-fundamental/discovery"
	"net/http"

	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type GRpcServerConfig struct {
	Name              string
	Address           string
	DisableLog        bool
	Logger            interfaces.Logger `json:"-" yaml:"-"`
	EnableCertVerify  bool
	CertInfo          *certutils.SecureOption
	MetaTransKeys     []string
	DiscoveryExConfig DiscoveryExConfig
}

type HttpServerConfig struct {
	Address           string
	Handler           http.Handler `json:"-" yaml:"-"`
	DiscoveryExConfig DiscoveryExConfig
}

type DiscoveryExConfig struct {
	Setter          discovery.Setter `json:"-" yaml:"-"`
	ExternalAddress string
	Meta            map[string]string
}

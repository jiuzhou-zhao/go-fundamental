package servicetoolset

import (
	"net/http"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/discovery"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/tracing"
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

	EnableTracing bool
	TracingConfig tracing.Config

	EnableGRpcWeb       bool
	GRpcWebAddress      string
	GRpcWebUseWebsocket bool
	GRpcWebPingInterval time.Duration
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

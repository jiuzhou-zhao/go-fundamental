package servicetoolset

import (
	"net/http"

	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type GRpcServerConfig struct {
	Address          string
	DisableLog       bool
	Logger           interfaces.Logger `json:"-" yaml:"-"`
	EnableCertVerify bool
	CertInfo         *certutils.SecureOption
	MetaTransKeys    []string
}

type HttpServerConfig struct {
	Address string
	Handler http.Handler `json:"-" yaml:"-"`
}

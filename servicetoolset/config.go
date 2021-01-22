package servicetoolset

import (
	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type GRpcConfig struct {
	Address          string
	DisableLog       bool
	Logger           interfaces.Logger `json:"-" yaml:"-"`
	EnableCertVerify bool
	CertInfo         *certutils.SecureOption
}

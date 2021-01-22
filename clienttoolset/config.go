package clienttoolset

import (
	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type GRpcClientConfig struct {
	Address       string
	SecureOption  certutils.SecureOption
	DisableLog    bool
	Logger        interfaces.Logger `json:"-" yaml:"-"`
	MetaTransKeys []string          `json:"-" yaml:"-"`
}

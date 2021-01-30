package clienttoolset

import (
	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"github.com/jiuzhou-zhao/go-fundamental/discovery"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type GRpcClientConfig struct {
	Address       string
	SecureOption  certutils.SecureOption
	DisableLog    bool
	Logger        interfaces.Logger `json:"-" yaml:"-"`
	MetaTransKeys []string          `json:"-" yaml:"-"`
}

type RegisterSchemasConfig struct {
	Getter  discovery.Getter  `json:"-" yaml:"-"`
	Logger  interfaces.Logger `json:"-" yaml:"-"`
	Schemas []string
}

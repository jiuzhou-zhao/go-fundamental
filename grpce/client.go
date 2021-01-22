package grpce

import (
	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"google.golang.org/grpc"
)

func DialGRpcServer(address string, secureOption *certutils.SecureOption, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOpts, err := DialOption(opts, secureOption)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(address, dialOpts...)
}

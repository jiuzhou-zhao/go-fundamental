package grpce

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/jiuzhou-zhao/go-fundamental/certutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func DialOption(opts []grpc.DialOption, secureOption *certutils.SecureOption) ([]grpc.DialOption, error) {
	if !secureOption.ServerWithTLS {
		return append(opts, grpc.WithInsecure()), nil
	}
	tlsConfig, err := generateClientTLSConfig(secureOption)
	if err != nil {
		return nil, err
	}
	return append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))), nil
}

func ServerOption(opts []grpc.ServerOption, secureOption *certutils.SecureOption) ([]grpc.ServerOption, error) {
	if !secureOption.ServerWithTLS {
		return opts, nil
	}
	tlsConfig, err := generateServerTLSConfig(secureOption)
	if err != nil {
		return opts, err
	}
	return append(opts, grpc.Creds(credentials.NewTLS(tlsConfig))), nil
}

func generateClientTLSConfig(secureOption *certutils.SecureOption) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	if secureOption == nil {
		return tlsConfig, nil
	}

	if !secureOption.VerifyServer {
		tlsConfig.InsecureSkipVerify = true
	}

	crt, err := tls.LoadX509KeyPair(secureOption.ClientCertFile, secureOption.ClientKeyFile)
	if err != nil {
		err = fmt.Errorf("generateClientTLSConfig LoadX509KeyPair failed: %v, %v, %v", err, secureOption.ClientCertFile, secureOption.ClientKeyFile)
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{crt}

	serverCertPool := x509.NewCertPool()
	for _, caFile := range secureOption.RootCAFiles {
		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			err = fmt.Errorf("generateClientTLSConfig load file failed: %v, %v", err, caFile)
			return nil, err
		}
		serverCertPool.AppendCertsFromPEM(caData)
	}

	for _, caFile := range secureOption.ServerCAFiles {
		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			err = fmt.Errorf("generateClientTLSConfig load file failed: %v, %v", err, caFile)
			return nil, err
		}
		serverCertPool.AppendCertsFromPEM(caData)
	}
	tlsConfig.RootCAs = serverCertPool

	if secureOption.ServerName != "" {
		tlsConfig.ServerName = secureOption.ServerName
	}

	return tlsConfig, nil
}

func generateServerTLSConfig(secureOption *certutils.SecureOption) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	if secureOption == nil {
		return tlsConfig, nil
	}

	if secureOption.VerifyClient {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
	}

	crt, err := tls.LoadX509KeyPair(secureOption.ServerCertFile, secureOption.ServerKeyFile)
	if err != nil {
		err = fmt.Errorf("generateServerTLSConfig LoadX509KeyPair failed: %v, %v, %v", err, secureOption.ServerCertFile, secureOption.ServerKeyFile)
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{crt}

	clientCertPool := x509.NewCertPool()
	for _, caFile := range secureOption.RootCAFiles {
		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			err = fmt.Errorf("generateServerTLSConfig load file failed: %v, %v", err, caFile)
			return nil, err
		}
		clientCertPool.AppendCertsFromPEM(caData)
	}

	for _, caFile := range secureOption.ClientCAFiles {
		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			err = fmt.Errorf("generateServerTLSConfig load file failed: %v, %v", err, caFile)
			return nil, err
		}
		clientCertPool.AppendCertsFromPEM(caData)
	}
	tlsConfig.ClientCAs = clientCertPool

	return tlsConfig, nil
}

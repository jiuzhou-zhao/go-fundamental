package grpce

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/jiuzhou-zhao/go-fundamental/discovery"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/iputils"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type BeforeServerStart func(server *grpc.Server)

type Server struct {
	serverName        string
	address           string
	logger            *loge.Logger
	beforeServerStart BeforeServerStart
	opts              []grpc.ServerOption

	setter          discovery.Setter
	externalAddress string
	meta            map[string]string
}

func NewServer(serverName string, address string, logger interfaces.Logger, beforeServerStart BeforeServerStart, opts []grpc.ServerOption) *Server {
	if serverName == "" {
		serverName = uuid.NewV4().String()
	}
	return &Server{
		serverName:        serverName,
		address:           address,
		logger:            loge.NewLogger(logger),
		beforeServerStart: beforeServerStart,
		opts:              opts,
	}
}

func (s *Server) EnableDiscovery(setter discovery.Setter, externalAddress string, meta map[string]string) {
	s.setter = setter
	s.externalAddress = externalAddress
	s.meta = meta
}

func (s *Server) getDiscoveryHostAndPort(ctx context.Context) (host string, port int, err error) {
	fnParse := func(address string) (host string, port int, err error) {
		if address == "" {
			err = errors.New("empty address")
			return
		}
		vs := strings.Split(address, ":")
		if len(vs) > 2 {
			err = errors.New("invalid address")
			return
		}
		if len(vs) == 2 {
			host = vs[0]
			var port64 int64
			port64, err = strconv.ParseInt(vs[1], 10, 64)
			if err != nil {
				return
			}
			port = int(port64)
		} else {
			host = address
		}
		return
	}
	host, port, err = fnParse(s.externalAddress)
	if err != nil {
		host = ""
		port = 0
		s.logger.Warnf(ctx, "parse external address %v failed: %v", s.externalAddress, err)
	}
	if host != "" && port > 0 {
		return
	}
	host2, port2, err := fnParse(s.address)
	if err != nil {
		return
	}
	if host == "" {
		host = host2
	}
	if port <= 0 {
		port = port2
	}
	if host == "" {
		ips, err := iputils.LocalIPv4s()
		if err == nil && len(ips) > 0 {
			host = ips[0]
		}
	}
	if host == "" || port < 0 {
		err = fmt.Errorf("invalid host port: %v,%v", host, port)
		return
	}
	return
}

func (s *Server) startDiscovery(ctx context.Context, server *grpc.Server) error {
	if s.setter == nil {
		return nil
	}
	host, port, err := s.getDiscoveryHostAndPort(ctx)
	if err != nil {
		return err
	}

	sis := server.GetServiceInfo()
	classV := ""
	for key := range sis {
		classV += "/" + key + ";"
	}
	if len(classV) > 0 {
		classV = classV[:len(classV)-1]
	}

	meta := map[string]string{discovery.MetaGRPCClass: classV}
	for k, v := range s.meta {
		meta[k] = v
	}
	return s.setter.Start([]*discovery.ServiceInfo{
		{
			Host:        host,
			Port:        port,
			ServiceName: discovery.BuildServerName(discovery.TypeGRpc, s.serverName, ""),
			Meta:        meta,
		},
	})
}

func (s *Server) stopAndWaitDiscovery() {
	if s.setter == nil {
		return
	}
	s.setter.Stop()
	s.setter.Wait()
}

func (s *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := grpc.NewServer(s.opts...)

	if s.beforeServerStart != nil {
		s.beforeServerStart(server)
	}

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return
	}
	s.logger.Infof(ctx, "grpc server listening on %v", s.address)

	go func() {
		err = s.startDiscovery(ctx, server)
		if err != nil {
			s.logger.Errorf(ctx, "discovery setter start failed: %v", err)
		}
		err = server.Serve(l)
		if err != nil {
			s.logger.Errorf(ctx, "grpc server serve error: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()

	s.logger.Infof(ctx, "grpc server shutting down")

	server.Stop()
	s.stopAndWaitDiscovery()

	return nil
}

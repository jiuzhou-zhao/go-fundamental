package toolset

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/clienttoolset"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type TestHelloWorld struct {
	helloworld.UnimplementedGreeterServer

	id            string
	helloWorldCli helloworld.GreeterClient
}

func (o *TestHelloWorld) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	respMsg := o.id
	if o.helloWorldCli != nil {
		resp, err := o.helloWorldCli.SayHello(ctx, &helloworld.HelloRequest{
			Name: o.id,
		})
		if err != nil {
			respMsg += err.Error()
		} else {
			respMsg = resp.Message
		}
	}
	return &helloworld.HelloReply{
		Message: fmt.Sprintf("Hi %v, I'm %v", req.Name, respMsg),
	}, nil
}

func TestGRpcServer1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	logger := &loge.ConsoleLogger{}

	go func() {
		conn, err := clienttoolset.DialGRpcServer(&clienttoolset.GRpcClientConfig{
			Address:    "127.0.0.1:9002",
			DisableLog: false,
			Logger:     logger,
		}, nil)
		assert.Nil(t, err)
		defer conn.Close()
		helloWorldCli := helloworld.NewGreeterClient(conn)

		serviceToolset := servicetoolset.NewServerToolset(ctx, logger)
		err = serviceToolset.CreateGRpcServer(&servicetoolset.GRpcServerConfig{
			Address:       ":9001",
			DisableLog:    false,
			MetaTransKeys: nil,
		}, nil, func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &TestHelloWorld{
				id:            "node1",
				helloWorldCli: helloWorldCli,
			})
		})

		assert.Nil(t, err)
		serviceToolset.Start()
		serviceToolset.Wait()
	}()

	go func() {
		serviceToolset := servicetoolset.NewServerToolset(ctx, &loge.ConsoleLogger{})
		err := serviceToolset.CreateGRpcServer(&servicetoolset.GRpcServerConfig{
			Address:       ":9002",
			DisableLog:    false,
			MetaTransKeys: nil,
		}, nil, func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &TestHelloWorld{id: "node2"})
		})
		assert.Nil(t, err)
		serviceToolset.Start()
		serviceToolset.Wait()
	}()

	time.Sleep(2 * time.Second)
	conn, err := clienttoolset.DialGRpcServer(&clienttoolset.GRpcClientConfig{
		Address:    "127.0.0.1:9001",
		DisableLog: false,
		Logger:     logger,
	}, nil)
	assert.Nil(t, err)
	defer conn.Close()
	helloWorldCli := helloworld.NewGreeterClient(conn)
	resp, err := helloWorldCli.SayHello(ctx, &helloworld.HelloRequest{Name: "cli"})
	assert.Nil(t, err)
	t.Log(resp.Message)
}

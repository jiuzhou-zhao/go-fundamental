package toolset

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jiuzhou-zhao/go-fundamental/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
	"testing"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/clienttoolset"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
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
	testGRpcServer1(t)
	time.Sleep(3 * time.Second)
}

func testGRpcServer1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	logger := &loge.ConsoleLogger{}

	tracingObj, _, err := tracing.NewGlobalTracer("httpTest", time.Second, "dev.env:6831")
	assert.Nil(t, err)

	go func() {
		conn, err := clienttoolset.DialGRpcServer(&clienttoolset.GRpcClientConfig{
			Address:       "127.0.0.1:9002",
			DisableLog:    false,
			Logger:        logger,
			EnableTracing: true,
		}, nil)
		assert.Nil(t, err)
		defer conn.Close()
		helloWorldCli := helloworld.NewGreeterClient(conn)

		serviceToolset := servicetoolset.NewServerToolset(ctx, logger)
		err = serviceToolset.CreateGRpcServer(&servicetoolset.GRpcServerConfig{
			Address:       ":9001",
			DisableLog:    false,
			MetaTransKeys: nil,
			EnableTracing: true,
		}, nil, func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &TestHelloWorld{
				id:            "node1",
				helloWorldCli: helloWorldCli,
			})
		})

		assert.Nil(t, err)
		_ = serviceToolset.Start()
		serviceToolset.Wait()
	}()

	go func() {
		serviceToolset := servicetoolset.NewServerToolset(ctx, &loge.ConsoleLogger{})
		err := serviceToolset.CreateGRpcServer(&servicetoolset.GRpcServerConfig{
			Address:       ":9002",
			DisableLog:    false,
			MetaTransKeys: nil,
			EnableTracing: true,
		}, nil, func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &TestHelloWorld{id: "node2"})
		})
		assert.Nil(t, err)
		_ = serviceToolset.Start()
		serviceToolset.Wait()
	}()

	go func() {
		serviceToolset := servicetoolset.NewServerToolset(ctx, &loge.ConsoleLogger{})
		r := mux.NewRouter()
		r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {

			conn, err := clienttoolset.DialGRpcServer(&clienttoolset.GRpcClientConfig{
				Address:       "127.0.0.1:9001",
				DisableLog:    false,
				Logger:        logger,
				EnableTracing: true,
			}, nil)
			assert.Nil(t, err)
			defer conn.Close()
			helloWorldCli := helloworld.NewGreeterClient(conn)

			resp, err := helloWorldCli.SayHello(request.Context(), &helloworld.HelloRequest{Name: "cli"})
			assert.Nil(t, err)
			t.Log(resp.Message)
		})

		err := serviceToolset.CreateHttpServer(&servicetoolset.HttpServerConfig{
			Address: ":8001",
			Handler: nethttp.Middleware(tracingObj, r),
		})
		assert.Nil(t, err)
		serviceToolset.Wait()
	}()

	time.Sleep(2 * time.Second)

	span := tracingObj.StartSpan("toplevel")
	span.SetTag(string(ext.Component), "client")

	req, err := http.NewRequest(
		"GET",
		"http://127.0.0.1:8001/test",
		nil,
	)
	assert.Nil(t, err)

	req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
	req, ht := nethttp.TraceRequest(tracingObj, req)

	client := &http.Client{Transport: &nethttp.Transport{}}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	_ = resp.Body.Close()

	ht.Finish()
	span.Finish()
}

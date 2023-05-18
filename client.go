package main

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/grpc-example/resolver/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const scheme = "ds"
const serviceName = "demo-service"

func main() {
	conn, err := grpc.Dial(scheme+"://"+serviceName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),                //禁用https
		grpc.WithResolvers(ds()),                                                //服务发现
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //负载均衡(新版本配置) round_robin/pick_first
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	name := "wang"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})

		if err != nil {
			w.Write([]byte(fmt.Sprintf("could not greet: %v", err)))
		} else {
			log.Println(r.Message)
			w.Write([]byte(r.Message + "\n"))
		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func ds() *manual.Resolver {
	m := manual.NewBuilderWithScheme(scheme)
	m.InitialState(resolver.State{Addresses: []resolver.Address{{Addr: "localhost:50001", ServerName: "demo11"}, {Addr: "localhost:50002"}}})
	m.BuildCallback = func(target resolver.Target, conn resolver.ClientConn, options resolver.BuildOptions) {
		log.Println(target, conn, options)
	}
	m.ResolveNowCallback = func(options resolver.ResolveNowOptions) {
		log.Println(options)
	}
	t := time.NewTimer(20 * time.Second)
	go func() {
		<-t.C
		m.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: "localhost:50003"}, {Addr: "localhost:50004"}}})
	}()
	return m
}

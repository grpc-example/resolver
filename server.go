package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/grpc-example/resolver/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	address string
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("request: ", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name + " port" + s.address}, nil
}

func main() {
	port := flag.Int("port", 50001, "端口号")
	flag.Parse()
	address := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpc.EnableTracing = true
	pb.RegisterGreeterServer(s, &server{address: address})
	log.Println("rpc服务已经开启", address)
	s.Serve(lis)
}

package grpc

import (
	"auth_service/grpc/pb"
	"fmt"
	"log"
	"net"
	"os"

	g "google.golang.org/grpc"
)

func GrpcListen() {

	listen := os.Getenv("AUTH_SERVICE_GRPC_LISTEN_IP") + ":" + os.Getenv("AUTH_SERVICE_GRPC_LISTEN_PORT")
	tcpListen, err := net.Listen("tcp", listen)

	if err != nil {
		log.Fatalf("failed to listen tcp on "+listen+": %v", err)
	}

	s := Server{}
	grpcServer := g.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	fmt.Println("grpc server is listening on " + listen)
	err = grpcServer.Serve(tcpListen)

	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

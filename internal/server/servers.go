package servers

import (
	gapi "github.com/razvan-bara/VUGO-API/api/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func LoadAuthGRPCServer(authService gapi.AuthServiceServer) {
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalln("Failed grpc tcp listen on users microservice.", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gapi.RegisterAuthServiceServer(grpcServer, authService)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln("Failed to  grpc server on users microservice.", err)
		}
	}()
}

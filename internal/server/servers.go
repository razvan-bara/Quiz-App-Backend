package servers

import (
	gapi "github.com/razvan-bara/VUGO-API/api/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func LoadAuthGRPCServer(authService gapi.AuthServiceServer) {
	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:4000")
		if err != nil {
			log.Fatalln("Failed grpc tcp listen on users microservice.", err)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		gapi.RegisterAuthServiceServer(grpcServer, authService)
		err = grpcServer.Serve(lis)
		grpcServer.Stop()
		if err != nil {
			log.Println("Failed to  grpc server on users microservice.", err)
		}
		log.Println("Grpc server starting")
	}()
}

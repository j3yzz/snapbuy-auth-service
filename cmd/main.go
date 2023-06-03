package main

import (
	"fmt"
	"github.com/j3yzz/snapbuy-auth-service/pkg/config"
	"github.com/j3yzz/snapbuy-auth-service/pkg/db"
	"github.com/j3yzz/snapbuy-auth-service/pkg/pb"
	"github.com/j3yzz/snapbuy-auth-service/pkg/services"
	"github.com/j3yzz/snapbuy-auth-service/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "snapbuy-auth",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth service on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

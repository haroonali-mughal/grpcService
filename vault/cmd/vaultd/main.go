package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vault"
	"vault/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	//cli "vault/cmd/vaultcli"
)

func main(){
	log.Println("1")
	var (
		httpAddr = flag.String("http", ":8080",
			"http listen address")
		gRPCAddr = flag.String("grpc", ":3000",
			"gRPC listen address")
	)
	log.Println("2")
	flag.Parse()
	ctx := context.Background()

	log.Println("before vault.NewService()")
	srv := vault.NewService()
	log.Printf("after vault.NewService() and return type is", reflect.TypeOf(srv),srv)

	errChan := make(chan error)

	go func() {
		log.Println("inside first go func() in main")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		log.Println("inside first go func() in main and before channel implementation")
		errChan <- fmt.Errorf("%s", <-c)
		log.Println("inside first go func() in main and after channel implementation")
	}()

	log.Println("before vault.MakeHashEndpoint(srv) in main")
	hashEndpoint := vault.MakeHashEndpoint(srv)
	log.Println("after vault.MakeHashEndpoint(srv) in main")

	log.Println("before vault.MakeValidateEndpoint(srv) in main")
	validateEndpoint := vault.MakeValidateEndpoint(srv)
	log.Println("after vault.MakeValidateEndpoint(srv) in main")

	endpoints := vault.Endpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}

	go func() {
		log.Println("http:", *httpAddr)
		handler := vault.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		log.Println("grpc:", *gRPCAddr)
		handler := vault.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	log.Fatalln(<-errChan)
	//cli.Main()

}

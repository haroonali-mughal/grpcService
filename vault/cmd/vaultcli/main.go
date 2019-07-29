package main

import (
	"flag"
	"fmt"
	grpc "google.golang.org/grpc"
	"log"
	"os"
	"reflect"
	"time"
	"context"
	"vault"
	grpcclient "vault/client/grpc"

	//"testing"

)

func main() {
	var (
		grpcAddr = flag.String("addr", ":3000",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()
	vaultService := grpcclient.New(conn)
	log.Println("i am printing :",reflect.TypeOf(vaultService))
	log.Println("hello inside another main function")
	log.Println("hello my name is khan")
	args := flag.Args()

	log.Println("these are command line args :", args)
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		log.Println("this is my password :",password)
		hash(ctx, vaultService, password)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		validate(ctx, vaultService, password, hash)
	default:
        log.Fatalln("unknown command", cmd)
 }
}

func hash(ctx context.Context, service vault.Service, password []string) {
	log.Println("inside hash password function")
	h, err := service.Hash(ctx, password[:])
	if err != nil {
		log.Fatalln("this is some error :",err.Error())
    }
    fmt.Println(h)
}

func validate(ctx context.Context, service vault.Service, password, hash []string) {
	valid, err := service.Validate(ctx, password[:], hash[:])
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !valid {
		fmt.Println("invalid")
		os.Exit(1)
	}
	fmt.Println("valid")
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}
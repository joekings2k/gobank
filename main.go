package main

import (
	"context"
	"database/sql"

	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joekings2k/gobank/api"
	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/gapi"
	"github.com/joekings2k/gobank/pb"
	"github.com/joekings2k/gobank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)


func main (){
	config,err :=util.LoadConfig(".")//path to config file
	if err != nil {
		log.Fatal("cannot log configs")
	}
	conn,err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil{
		log.Fatal("cannot connect to db:",err)
	}
	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
	// runGinServer(config,store)
}


func runGrpcServer(config util.Config, store db.Store){
	
	server ,err := gapi.NewServer(config ,store)
	if err !=nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGobankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store){
	
	server ,err := gapi.NewServer(config ,store)
	if err !=nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption:= runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterGobankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(config util.Config, store db.Store){
	
	server ,err := api.NewServer(config ,store)
	if err !=nil {
		log.Fatal("cannot create server:", err)
	}
	err= server.Start(config.HTTPServerAddress)
	if err !=nil {
		log.Fatal("cannot start server:", err)
	}
}
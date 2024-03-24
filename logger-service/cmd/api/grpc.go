package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dineshkuncham/logger-service/data"
	"github.com/dineshkuncham/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}
	res := &logs.LogResponse{
		Result: "Logged!",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listem for gRPC:%v", err)
	}
	sv := grpc.NewServer()

	logs.RegisterLogServiceServer(sv, &LogServer{
		Models: app.Models,
	})
	log.Printf("gRPC server started on port %s", grpcPort)
	if err := sv.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for gRPC:%v", err)
	}
}

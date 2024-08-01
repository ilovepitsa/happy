package app

import (
	"net"

	"github.com/ilovepitsa/happy/auth/api/sessions"
	grpchandler "github.com/ilovepitsa/happy/auth/internal/handlers/grpcHandler"
	kv "github.com/ilovepitsa/happy/auth/internal/repo/KV"
	"github.com/ilovepitsa/happy/auth/internal/service"
	"github.com/ilovepitsa/happy/auth/pkg/config"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

func Run(configPath string) error {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}
	SetLogrusParams(cfg)

	log.Info("Initializing repo service.....")
	r, err := kv.New(cfg)
	if err != nil {
		return err
	}
	defer r.Close()

	log.Info("Initializing session service.....")
	sessManager := service.NewKVSessManager(r)

	log.Info("Initializing session grpc.....")
	grpcSess := grpchandler.NewSessionHandler(sessManager)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	log.Info("Starting listing tcp ", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	if err != nil {
		return err
	}

	sessions.RegisterAuthCheckerServer(grpcServer, grpcSess)

	err = grpcServer.Serve(lis)
	return err
}

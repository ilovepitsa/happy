package app

import (
	"net"

	"github.com/ilovepitsa/happy/notify/pkg/config"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

func Run(configPath string) error {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}
	SetLogrusParams(cfg)

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	log.Info("Starting listing tcp ", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	if err != nil {
		return err
	}

	err = grpcServer.Serve(lis)
	return err
}

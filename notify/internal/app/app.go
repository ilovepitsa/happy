package app

import (
	"net"

	"github.com/ilovepitsa/happy/notify/api/notifier"
	grpchandler "github.com/ilovepitsa/happy/notify/internal/handlers/grpcHandler"
	"github.com/ilovepitsa/happy/notify/internal/repo"
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

	log.Info("Initializing repo...")
	r, err := repo.NewRepo(cfg.P)
	if err != nil {
		log.Fatal("error init repo ", err)
	}

	grpcServer := grpc.NewServer(opts...)
	log.Info("Starting listing tcp ", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.Net.Host, cfg.Net.Port))
	if err != nil {
		return err
	}

	notifierGrpc := grpchandler.NewNotificationServer(r.NotifRepo)

	notifier.RegisterNotifierServer(grpcServer, notifierGrpc)

	err = grpcServer.Serve(lis)
	return err
}

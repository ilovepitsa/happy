package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ilovepitsa/happy/backend/pkg/config"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) error {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}
	SetLogrusParams(cfg)

	log.Info("Initializing repo service.....")
	router := gin.Default()

	services

	return err
}

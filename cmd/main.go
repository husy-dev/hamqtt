package main

import (
	"context"

	"github.com/husy-dev/hamqtt/internal/brokers"
	"github.com/husy-dev/hamqtt/internal/configs"
	"github.com/husy-dev/hamqtt/internal/logger"
	"go.uber.org/zap"
)
var log = logger.Get()

func main() {
	config, err := configs.LoadConf()
	if err != nil {
		log.Fatal("configure broker config error", zap.Error(err))
	}
	ctx:=context.Background()
	b := brokers.NewBroker(config,ctx)
	b.Start()
}

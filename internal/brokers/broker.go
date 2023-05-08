package brokers

import (
	"context"
	"net"

	"github.com/husy-dev/hamqtt/internal/clients"
	"github.com/husy-dev/hamqtt/internal/configs"
	"github.com/husy-dev/hamqtt/internal/logger"
	"github.com/husy-dev/hamqtt/internal/topics"
	"go.uber.org/zap"
)
var log = logger.Get()

type Broker struct {
	clientMgr *clients.ClientMgr
	config  *configs.Config
	topicMgr *topics.TopicMgr
	ctx context.Context
}

func NewBroker(c *configs.Config,ctx context.Context)*Broker{
	return &Broker{
		config: c,
		ctx:ctx,
		clientMgr: clients.NewClientMgr(),
		topicMgr: topics.NewTopicMgr(),
	}
}

func(b *Broker) Start(){
	hp := b.config.Host + ":" + b.config.Port
	l, err := net.Listen("tcp", hp)
	if err!=nil{
		log.Fatal("broker listening failed!",zap.Error(err))
	}
	log.Info("Start Listening client on ", zap.String("hp", hp))
	for {
		select {
		case <-b.ctx.Done():
			//TODO: 做一些收尾工作：关闭所有连接的clients
			b.clientMgr.Close()
			return
		default:
			conn, err := l.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					//TODO: 做一个等待
				} else {
					log.Error("Accept error", zap.Error(err))
				}
				continue
			}
			c:=b.clientMgr.GetClient(conn,b.ctx)
			go c.ReadLoop()
		}
	}
}

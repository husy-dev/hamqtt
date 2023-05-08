package clients

import (
	"context"
	"io"
	"net"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/husy-dev/hamqtt/internal/logger"
	"go.uber.org/zap"
)
var log = logger.Get()

type Client struct {
	clientID string
	conn     net.Conn
	ctx  context.Context
}


type InflightStatus uint8
const (
	Publish InflightStatus = 0
	Pubrel  InflightStatus = 1
)

func (c *Client)Init(conn net.Conn,ctx context.Context){
	c.conn = conn
	c.ctx = ctx
}

func (c *Client) ReadLoop() {
	// keepAlive := time.Second * time.Duration(c.info.keepalive)
	// timeOut := keepAlive + (keepAlive / 2)

	for {
		select {
		case <-c.ctx.Done():
			//TODO: 做一些收尾工作
			c.close()
			return
		default:
			packet, err := packets.ReadPacket(c.conn)
			//TODO：读取失败，EOF的处理存疑
			if err != nil && err!=io.EOF{
				log.Error("read packet error: ", zap.Error(err), zap.String("ClientID", c.clientID))
				c.close()
				return
			}
			c.handleMessage(packet)
		}
	}

}

func(c *Client)handleMessage(p packets.ControlPacket){
	switch p.(type) {
	case *packets.ConnackPacket:
	case *packets.ConnectPacket:
	case *packets.PublishPacket:
		c.ProcessPublish(p.(*packets.PublishPacket))
	case *packets.PubackPacket:
	case *packets.PubrecPacket:
	case *packets.PubrelPacket:
	case *packets.PubcompPacket:
	case *packets.SubscribePacket:
		c.ProcessSubscribe( p.(*packets.SubscribePacket))
	case *packets.SubackPacket:
	case *packets.UnsubscribePacket:
		c.ProcessUnSubscribe(p.(*packets.UnsubscribePacket))
	case *packets.UnsubackPacket:
	case *packets.PingreqPacket:
		c.ProcessPing()
	case *packets.PingrespPacket:
	case *packets.DisconnectPacket:
		c.close()
	default:
		log.Info("Recv Unknow message.......", zap.String(logger.CLIENTID, c.clientID))
	}
}


func(c *Client)ProcessPublish(p *packets.PublishPacket){
	// topic := p.TopicName
	// if !c.broker.CheckTopicAuth(PUB, c.clientID, c.info.username, c.info.remoteIP, topic) {
	// 	log.Error("Pub Topics Auth failed, ", zap.String("topic", topic), zap.String("ClientID", c.info.clientID))
	// 	return
	// }

	switch p.Qos {
	case QosAtMostOnce:
		c.ProcessPublishMessage(p)
	case QosAtLeastOnce:
		puback := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
		puback.MessageID = p.MessageID
		if err := c.WriterPacket(puback); err != nil {
			log.Error("send puback error, ", zap.Error(err), zap.String("ClientID", c.clientID))
			return
		}
		c.ProcessPublishMessage(p)
	case QosExactlyOnce:
		if err := c.registerPublishPacketId(p.MessageID); err != nil {
			return
		} else {
			pubrec := packets.NewControlPacket(packets.Pubrec).(*packets.PubrecPacket)
			pubrec.MessageID = p.MessageID
			if err := c.WriterPacket(pubrec); err != nil {
				log.Error("send pubrec error, ", zap.Error(err), zap.String("ClientID", c.clientID))
				return
			}
			c.ProcessPublishMessage(p)
		}
		return
	default:
		log.Error("publish with unknown qos", zap.String("ClientID", c.clientID))
		return
	}
}

//TODO：考虑资源回收、消息未发送的记录
func(c *Client)close() error{
	err:=c.conn.Close()
	if(err!=nil){
		log.Error("close client failed",zap.String(logger.CLIENTID,c.clientID),zap.Error(err))
	}
	return err
}
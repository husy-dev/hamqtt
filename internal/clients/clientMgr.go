package clients

import (
	"context"
	"net"
	"sync"
)

type ClientMgr struct{
	clients sync.Pool
	clientsMap sync.Map
}

func NewClientMgr() *ClientMgr{
	return &ClientMgr{
		clients: sync.Pool{
			New: func()interface{}{
				return &Client{}
			},
		},
	}
}

func(cm *ClientMgr)GetClient(conn net.Conn,ctx context.Context)*Client{
	c:= cm.clients.Get().(*Client)
	c.Init(conn,ctx)
	return c
}


func(cm *ClientMgr)Close(){
	//TODO：关闭所有client，每一个client自己做好自己的收尾工作
	cm.clientsMap.Range(func(key, value any) bool {
		value.(*Client).close()
		return true
	})
}
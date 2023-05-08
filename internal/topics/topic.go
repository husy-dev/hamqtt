package topics

import "github.com/husy-dev/hamqtt/internal/clients"

// TopicsProvider
type Topic struct {
	topic string
	subs []*clients.Client
}

func NewTopic(topic string)*Topic{
	return &Topic{
		topic: topic,
	}
}

func (t *Topic) GetSubs()[]*clients.Client{
	return t.subs
}

func(t *Topic)AddSub(c *clients.Client){
	t.subs= append(t.subs, c)
	return
}

func(t *Topic)CheckTopicAuth(c *clients.Client){
	
}

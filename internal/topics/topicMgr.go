package topics

import "github.com/husy-dev/hamqtt/internal/clients"

type TopicMgr struct{
	topics map[string]*Topic
}

func NewTopicMgr() (*TopicMgr) {
	return &TopicMgr{
		topics:make(map[string]*Topic),
	}
}

func(t *TopicMgr)Subscribe(topic string, c *clients.Client){
	if old,exist:= t.topics[topic];exist{
		old.AddSub(c)
	}else{
		t.topics[topic]=NewTopic(topic)
		t.topics[topic].AddSub(c)
	}
}

func(t *TopicMgr)Unsubscribe(topic string, c *clients.Client){
	
}

func(t *TopicMgr)GetSubs(topic string)[]*clients.Client{
	return t.topics[topic].GetSubs()
}
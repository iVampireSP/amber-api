package stream

import (
	"context"
	"fmt"
	"rag-new/internal/base/conf"
	logger2 "rag-new/internal/base/logger"
	"testing"
)

func NewStream() *Service {
	logger := logger2.NewZapLogger()
	config := conf.ProviderConfig(logger)
	//database := orm.NewGORM(config, logger)
	//dao := dao2.NewQuery(database)
	//milvus := milvus2.NewMilvus(config)

	return NewService(config)
}

//	func TestListen(t *testing.T) {
//		stream := NewStream()
//
//		var topic = "test"
//
//		err := stream.Listen(topic, callback)
//		if err != nil {
//			t.Fatal(err)
//		}
//
// }
func callback(data []byte) {
	fmt.Println(string(data))
}

//	func TestPublish(t *testing.T) {
//		stream := NewStream()
//
//		var topic = "test"
//
//		for {
//			err := stream.Publish(topic, []byte("Hello"))
//			if err != nil {
//				t.Fatal(err)
//				return
//			}
//		}
//	}
func TestProducer(t *testing.T) {
	stream := NewStream()

	var ctx = context.Background()

	var m = 100000

	for i := 0; i < m; i++ {
		//producer := stream.Producer(topic)
		err := stream.SendMessage(ctx, nil, []byte("Hello"))
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}

func TestConsumer(t *testing.T) {
	stream := NewStream()

	var topic = "test"

	var ctx = context.Background()

	stream.ReadMessage(ctx, topic, "test")
}

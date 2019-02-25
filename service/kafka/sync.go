package kafka

import (
	"log"
	"fmt"
	"time"
	"github.com/Shopify/sarama"
	"os"
)



//同步消息模式
func SyncProducer()  {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(Address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	topic := Topic
	srcValue := "sync: this is a message. index=%d"
	for i:=0; i<10; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic:topic,
			Value:sarama.ByteEncoder(value),
//			Key:sarama.ByteEncoder("aaaa"),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		}else {
			fmt.Fprintf(os.Stdout, value + "发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2*time.Second)
	}
}
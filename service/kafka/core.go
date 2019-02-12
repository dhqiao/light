package kafka

import (
	"time"
	"github.com/Shopify/sarama"
	"fmt"
	"strconv"
	"github.com/bsm/sarama-cluster"
)

var Address = []string{"127.0.0.1:9092","127.0.0.1:9091","127.0.0.1:9094"}

var Topic = "test"



// 异步模式
func SaramaProducer() {
	producer, e := NewAsyncProducer()
	if e != nil {
		fmt.Println(e)
		return
	}
	defer producer.AsyncClose()


	//循环判断哪个通道发送过来数据. 监控发送数据的返回是否有问题，如果出现问题，可以做一些事情
	fmt.Println("start goroutine")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				fmt.Println("$$$$ offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(producer)

	var value string
	for i := 0; ; i++ {
		if i == 50 {
			return
		}
		time.Sleep(500 * time.Millisecond)
		value = strconv.Itoa(i) + " this is a message 0606 " + time.Now().Format("15:04:05")

		// 发送的消息,主题。
		// 注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系。
		msg := &sarama.ProducerMessage{
			Topic: Topic,
		}

		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(value)

		//使用通道发送
		producer.Input() <- msg
	}

}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_1_0_0

	fmt.Println("start make producer")
	//使用配置,新建一个异步生产者
	return sarama.NewAsyncProducer(Address, config)
}

func NewConsumer(groupId string) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	topics := []string{Topic}
	brokers := Address
	// init consumer
	return cluster.NewConsumer(brokers, groupId, topics, config)
}
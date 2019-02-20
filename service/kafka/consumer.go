package kafka

import (
	"sync"
	"log"
	"os"
	"os/signal"
	"fmt"
)

// 支持brokers cluster的消费者
func ClusterConsumer(wg *sync.WaitGroup, groupId string) {
	defer wg.Done()
	// init consumer
	consumer, err := NewConsumer(groupId)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			fmt.Println("...........................")
			if ok {
				fmt.Println("------------------------")
				fmt.Fprintf(os.Stdout, "...... %s:%s:%s/%d/%d\t%s\t%s\n", msg.Key, groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
				fmt.Println("------------------------")
			}
		case <-signals:
			fmt.Println("0000000000000000000000000000000")
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)

}
package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"sync/atomic"
	"time"
)

var (
	kafkaSender *KafkaSender
)

type KafkaSender struct {
	client sarama.SyncProducer
	lineChan chan Message
}
type Message struct {
	t string
	msg string
}

// 初始化kafka
func NewKafkaSender()(kafka *KafkaSender,err error){
	kafka = &KafkaSender{
		lineChan:make(chan Message,10000),
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client,err := sarama.NewSyncProducer([]string{"tdh2:9092","tdh3:9092","tdh4:9092","tdh5:9092","tdh6:9092","tdh7:9092","tdh8:9092","tdh9:9092","tdh10:9092","tdh11:9092","tdh12:9092","tdh14:9092","tdh13:9092"},config)
	if err != nil{
		logs.Error("init kafka client failed,err:%v\n",err)
		return
	}
	kafka.client = client
	stati := statiBenchmarkProducerSnapshot{}
	snapshots := produceSnapshots{cur: &stati}
	// 根据CPU来协商启动协程数
	for i := 0; i < 24; i++ {
		go func() {
			waitgroup.Add(1)
			kafka.sendToKafka("test_op",&stati)
			kafka.sendToKafka("test_zt",&stati)
			waitgroup.Done()
		}()
	}
	// snapshot
	go func() {
		waitgroup.Add(1)
		defer waitgroup.Done()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				snapshots.TakeSnapshot()
			}
		}
	}()
	// print statistic
	go func() {
		waitgroup.Add(1)
		defer waitgroup.Done()
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				snapshots.TrintStati()
			}
		}
	}()

	return
}

func InitKafka()(err error){
	kafkaSender,err = NewKafkaSender()
	return
}

func (k *KafkaSender) sendToKafka(topic string,stati *statiBenchmarkProducerSnapshot){
	//从channel中读取日志内容放到kafka消息队列中
	count:=0
	now := time.Now()
	for v := range k.lineChan{
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.t
		msg.Value = sarama.StringEncoder(v.msg)
		p,offset,err := k.client.SendMessage(msg)
		if err==nil{
			count++;
			atomic.AddInt64(&stati.receiveResponseSuccessCount, 1)
			atomic.AddInt64(&stati.sendRequestSuccessCount, 1)
			currentRT := int64(time.Since(now) / time.Millisecond)
			atomic.AddInt64(&stati.sendMessageSuccessTimeTotal, currentRT)
			prevRT := atomic.LoadInt64(&stati.sendMessageMaxRT)
			for currentRT > prevRT {
				if atomic.CompareAndSwapInt64(&stati.sendMessageMaxRT, prevRT, currentRT) {
					break
				}
				prevRT = atomic.LoadInt64(&stati.sendMessageMaxRT)
			}
		}
		if count % 10000 == 0{
			logs.Info(" offset %d  partition %d ",offset,p)
		}
		if err != nil{
			logs.Error("send message to kafka failed,err:%v",err)
		}
	}
}

func (k *KafkaSender) addMessage(line string,topic string)(err error){
	//我们通过tailf读取的日志文件内容先放到channel里面
	k.lineChan <- Message{t:topic,msg:line}
	return
}
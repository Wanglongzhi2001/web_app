package RabbitMQ

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"sync"
	"web_app/logic"
	"web_app/models"
)

// url 格式 amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://testuser:testuser@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 连接信息
	Mqurl string

	sync.Mutex
}

// 构造函数创建RabbitMQ实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ { // 根据这三个参数的不同组合实现不同的模式
	rabbitmq := RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	// 创建RabbitMQ连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "获取channel失败")
	return &rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.conn.Close()
	r.channel.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Simple模式step1:在Simple模式下创建RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "") // exchange使用默认的(direct),key为空(不同模式传入不同参数)
}

// Simple模式step2:在Simple模式下生产代码
func (r *RabbitMQ) PublishSimple(msg string) {

	// 1.声明队列
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除(当最后一个订阅者取消订阅时)
		false,       // 是否具有排他性
		false,       // 是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 如果为true，根据exchange类型和routekey规则，如果找不到符合条件的队列那么会把发送的消息返回给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

}

// 消费消息
func (r *RabbitMQ) ConsumeSimple() {
	var err error
	// 1.声明队列
	_, err = r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除(当最后一个订阅者取消订阅时)
		false,       // 是否具有排他性
		false,       // 是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列是否阻塞
		false,
		// 其他参数
		nil)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 启用协程处理消息(因为我们生产消息是异步的)
	post := new(models.Post)
	go func() {
		for d := range msgs {
			// 实现我们的逻辑函数
			if err = json.Unmarshal([]byte(d.Body), post); err != nil {
				zap.L().Error("反序列化消息失败！", zap.Error(err))
			}
			fmt.Println(post.ID)
			err = logic.CreatePost(post)
			if err != nil {
				zap.L().Error("logic.CreatePost failed", zap.Error(err))
				return
			}
			fmt.Println("已经收到了数据，帖子ID为：", post.ID)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	// 使之阻塞
	<-forever
}

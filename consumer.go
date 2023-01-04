package main

import (
	"web_app/models"
	RabbitMQ "web_app/rabbitmq"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple(models.QueueName) // 名字必须一样
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeSimple()
}

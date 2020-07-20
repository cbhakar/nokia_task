package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
	"nokia_task/service"
)

var (
	url       = "amqp://ssvvtrpr:aW-jF853zRnQ8LHajrYfvAwi04bknIZn@lionfish.rmq.cloudamqp.com/ssvvtrpr"
	queueName = "test"
	exchange  = "nokia"
)

func RabbitMqinit() {
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	channel, err := connection.Channel()
	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	log.Println("rabbitmq channel created")

	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}
	log.Println("queue created with name : ", queueName)

	err = channel.QueueBind(queueName, "#", exchange, false, nil)
	if err != nil {
		panic("error binding queue to routing key" + err.Error())
	}
	log.Println("queue binded with exchange : ", exchange)

	msgs, err := channel.Consume("test", "", false, false, false, false, nil)
	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	go CheckMsg(msgs)
}

func CheckMsg(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		log.Println("task recieved on queue")
		err := service.ReloadDataToRedis()
		if err != nil {
			log.Println("error reloading user data o redis")
		}
		err = msg.Ack(false)
		if err != nil {
			log.Println("error acknowledging  rabbitmq msg")
		}
		log.Println("reload task completed")
	}
}

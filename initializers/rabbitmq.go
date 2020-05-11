package initializers

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/oldfritter/sneaker-go/utils"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
)

type Amqp struct {
	Connect struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Vhost    string `yaml:"vhost"`
	} `yaml:"connect"`

	Exchange map[string]map[string]string `yaml:"exchange"`
}

var (
	AmqpGlobalConfig Amqp
	RabbitMqConnect  utils.RabbitMqConnect
)

func InitializeAmqpConfig() {
	path_str, _ := filepath.Abs("config/amqp.yml")
	content, err := ioutil.ReadFile(path_str)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal(content, &AmqpGlobalConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	InitializeAmqpConnection()
}

func InitializeAmqpConnection() {
	var err error
	conn, err := amqp.Dial("amqp://" + AmqpGlobalConfig.Connect.Username + ":" + AmqpGlobalConfig.Connect.Password + "@" + AmqpGlobalConfig.Connect.Host + ":" + AmqpGlobalConfig.Connect.Port + "/" + AmqpGlobalConfig.Connect.Vhost)
	RabbitMqConnect = utils.RabbitMqConnect{conn}
	if err != nil {
		time.Sleep(5000)
		InitializeAmqpConnection()
		return
	}
	go func() {
		<-RabbitMqConnect.NotifyClose(make(chan *amqp.Error))
		InitializeAmqpConnection()
	}()
}

func CloseAmqpConnection() {
	RabbitMqConnect.Close()
}

func GetRabbitMqConnect() utils.RabbitMqConnect {
	return RabbitMqConnect
}

func PublishMessageWithRouteKey(exchange, routeKey, contentType string, message *[]byte, arguments amqp.Table, deliveryMode uint8) error {
	channel, err := RabbitMqConnect.Channel()
	defer channel.Close()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	if err = channel.Publish(
		exchange, // publish to an exchange
		routeKey, // routing to 0 or more queues
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     contentType,
			ContentEncoding: "",
			Body:            *message,
			DeliveryMode:    deliveryMode, // amqp.Persistent, amqp.Transient // 1=non-persistent, 2=persistent
			Priority:        0,            // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Queue Publish: %s", err)
	}
	return nil
}

func PublishMessageToQueue(queue, contentType string, message *[]byte, arguments amqp.Table, deliveryMode uint8) error {
	channel, err := RabbitMqConnect.Channel()
	defer channel.Close()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	if err = channel.Publish(
		"",    // publish to an exchange
		queue, // routing to 0 or more queues
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     contentType,
			ContentEncoding: "",
			Body:            *message,
			DeliveryMode:    deliveryMode, // amqp.Persistent, amqp.Transient // 1=non-persistent, 2=persistent
			Priority:        0,            // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Queue Publish: %s", err)
	}
	return nil
}

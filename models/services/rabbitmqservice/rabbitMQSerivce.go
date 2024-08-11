package rabbitmqservice

import (
	"encoding/json"
	"os"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/rabbitmq/amqp091-go"
)

type QueueMessageType string 

const (
	POST QueueMessageType = "POST"
)
type RabbitMQPostMessage struct {
    Type QueueMessageType
    Poster string
}

type RabbitMQService struct {
    Conn *amqp091.Connection
    Chan *amqp091.Channel
    QName string
}

type IRabbitMQService interface {
    Connect() 
    Enqueue(v any) error
}
func(rmq *RabbitMQService) Enqueue(v any) error {
    bytes, err := json.Marshal(v)
    if err != nil {
        return customerrors.WrapBasicError(err)
    }
    return rmq.Chan.Publish(
        "",
        rmq.QName,
        false,
        false,
        amqp091.Publishing{
            ContentType: "application/json",
            Body: bytes,
        },
    )
}

func(rmq *RabbitMQService) Connect() {
    connectionString := os.Getenv("RABBIT_MQ_CONNECTION_STRING")

    conn, err := amqp091.Dial(connectionString)

    if err != nil {
        panic(err.Error())
    }

    rmq.Conn = conn

    ch, err := conn.Channel()

    if err != nil {
        panic(err.Error())
    }

    rmq.Chan = ch
    rmq.QName = "emailQueue"

    _, err = ch.QueueDeclare(
      rmq.QName, 
      false,   
      false,   
      false,  
      false, 
      nil,  
    )
}


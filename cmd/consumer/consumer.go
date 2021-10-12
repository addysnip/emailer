package consumer

import (
	"fmt"

	"addysnip.dev/api/pkg/utils"
	"addysnip.dev/emailer/cmd/migrate"
	"addysnip.dev/emailer/pkg/logger"
	"addysnip.dev/emailer/pkg/mailer"
	"github.com/streadway/amqp"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "consumer",
		Usage:  "Run message consumer",
		Action: Run,
	}
}

var log = logger.Category("cmd/consumer")

func Run(c *cli.Context) error {
	err := migrate.Run(c)
	if err != nil {
		return err
	}

	log.Info("Configuring RabbitMQ Consumer")
	connstring := fmt.Sprintf("amqp://%s:%s@%s:%s",
		utils.Getenv("RABBITMQ_USER", "guest"),
		utils.Getenv("RABBITMQ_PASSWORD", "guest"),
		utils.Getenv("RABBITMQ_HOST", "localhost"),
		utils.Getenv("RABBITMQ_PORT", "5672"),
	)
	conn, err := amqp.Dial(connstring)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		utils.Getenv("RABBITMQ_QUEUE", "emailer"), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	log.Info("Channel and Queue Established")
	defer conn.Close()
	defer ch.Close()

	log.Info("Building consumer")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() error {
		log.Info("Consumer started")
		for d := range msgs {
			log.Debug("Received a message: %s", d.Body)

			err := mailer.Handle(string(d.Body))
			if err != nil {
				log.Error("Error handling message: %s, re-queueing", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
		return nil
	}()

	log.Info("Consumer is ready, entering goroutine...")
	<-forever

	return nil
}

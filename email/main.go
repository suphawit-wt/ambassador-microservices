package main

import (
	"ambassador/email/models"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gobuffalo/envy"
)

func main() {
	envy.Load()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_PROTOCOL"),
		"sasl.username":     os.Getenv("KAFKA_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_PASSWORD"),
		"sasl.mechanism":    os.Getenv("KAFKA_MECHANISM"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{os.Getenv("KAFKA_TOPIC")}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			message := models.EmailKafkaMessage{}

			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				panic(err)
			}

			smtpHost := os.Getenv("SMTP_HOST")
			smtpPort := os.Getenv("SMTP_PORT")

			smtpAuth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), smtpHost)

			ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", message.AmbassadorRevenue, message.Code))

			smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, "no-reply@email.com", []string{message.AmbassadorEmail}, ambassadorMessage)

			adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", message.Id, message.AdminRevenue))

			smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()
}

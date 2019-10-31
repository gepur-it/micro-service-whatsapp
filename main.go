package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/zbindenren/logrus_mail"
	"log"
	"net/http"
	"os"
)

var AMQPConnection *amqp.Connection
var AMQPChannel *amqp.Channel

var logger = logrus.New()

func init()  {
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")

	port, err := strToInt("LOGTOEMAIL_SMTP_PORT")
	failOnError(err, "Error read smtp port from env")

	hook, err := logrus_mail.NewMailAuthHook(
		os.Getenv("LOGTOEMAIL_APP_NAME"),
		os.Getenv("LOGTOEMAIL_SMTP_HOST"),
		port,
		os.Getenv("LOGTOEMAIL_SMTP_FROM"),
		os.Getenv("LOGTOEMAIL_SMTP_TO"),
		os.Getenv("LOGTOEMAIL_SMTP_USERNAME"),
		os.Getenv("LOGTOEMAIL_SMTP_PASSWORD"),
	)

	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})

	logger.Hooks.Add(hook)

	cs := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_ERP_LOGIN"),
		os.Getenv("RABBITMQ_ERP_PASS"),
		os.Getenv("RABBITMQ_ERP_HOST"),
		os.Getenv("RABBITMQ_ERP_PORT"),
		os.Getenv("RABBITMQ_ERP_VHOST"))

	connection, err := amqp.Dial(cs)
	failOnError(err, "Failed to connect to RabbitMQ")
	AMQPConnection = connection


	channel, err := AMQPConnection.Channel()
	failOnError(err, "Failed to open a channel")
	AMQPChannel = channel
}

func main() {
	http.HandleFunc("/", receiver)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT")), nil))
}

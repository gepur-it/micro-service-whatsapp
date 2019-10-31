package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strings"
	"time"
)

type Command struct {
	Id      string                 `json:"id"`
	Payload map[string]interface{} `json:"payload"`
}

func receiver(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")

	if id == "" {
		log.Printf("Received a callback with empty id: %s", id)
	} else {
		log.Printf("Received a callback id: %s", id)

		callbackRequest := map[string]interface{}{}

		err := json.NewDecoder(r.Body).Decode(&callbackRequest)
		failOnError(err, "Can`t decode webHook callBack")

		command := Command{id, callbackRequest}

		callbackRequestJson, err := json.Marshal(command)
		failOnError(err, "Can`t serialise webHook response")

		log.Printf("Received a callback: %s", callbackRequestJson)

		name := fmt.Sprintf("whats_app_to_erp")

		err = AMQPChannel.Publish(
			"",
			name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Transient,
				ContentType:  "application/json",
				Body:         callbackRequestJson,
				Timestamp:    time.Now(),
			})

		failOnError(err, "Failed to publish a message")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(callbackRequestJson)
		failOnError(err, "Failed to write API answer")
	}
}

package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

type Message struct {
	Type     string
	BannerId int
	SlotId   int
	GroupId  int
}

var ErrUserNotFoundInProducer = errors.New("user not found")

func SendKafkaMessage(producer sarama.SyncProducer, message Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(message.BannerId)),
		Value: sarama.StringEncoder(messageJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func SendMessageHandler(producer sarama.SyncProducer, message Message) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := SendKafkaMessage(producer, message)
		if errors.Is(err, ErrUserNotFoundInProducer) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

package kafka

import (
	"github.com/IBM/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
)

func NewCloudEventsKafkaConsumer(sc *SaramaConfig) (CloudEventsReceiver, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	receiver, err := kafka_sarama.NewConsumer([]string{sc.Host}, saramaConfig, sc.GroupId, sc.Topic)
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(receiver)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c, nil
}

func NewCloudEventsKafkaSender(sc *SaramaConfig) (CloudEventsSender, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	sender, err := kafka_sarama.NewSender([]string{sc.Host}, saramaConfig, sc.Topic)
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c, nil
}

type CloudEventsReceiver cloudevents.Client

type CloudEventsSender cloudevents.Client

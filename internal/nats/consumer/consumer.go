package consumer

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"img_downloader/internal/config"
	"log/slog"
	"strconv"
)

type Consumer struct {
	log   *slog.Logger
	conn  *nats.Conn
	topic string
}

func New(natsCfg *config.NatsConfig, topic string, log *slog.Logger) *Consumer {
	natsURL := fmt.Sprintf("nats://%s:%s@%s:%d", natsCfg.User, natsCfg.Password, natsCfg.Host, natsCfg.Port)
	conn, err := nats.Connect(natsURL)

	if err != nil {
		log.Error("Failed to create nats consumer", err)
		panic(err)
	}

	log.Info("Created NATS consumer on", slog.String("Port", strconv.Itoa(natsCfg.Port)))

	return &Consumer{
		conn:  conn,
		topic: topic,
		log:   log,
	}
}

func (c *Consumer) Start() {
	log := c.log.With(
		slog.String("topic", c.topic),
	)

	_, err := c.conn.Subscribe(c.topic, func(msg *nats.Msg) {
		url := string(msg.Data)
		log.Info("Received a message", slog.String("url", url))
	})

	if err != nil {
		log.Error(fmt.Sprintf("Failed to subscribe: %v", err))
		return
	}
}

package consumer

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"img_downloader/internal/config"
	"log/slog"
	"strconv"
)

type Handler interface {
	Process(ctx context.Context, msg *nats.Msg) error
}

type Consumer struct {
	log      *slog.Logger
	conn     *nats.Conn
	handlers map[string]Handler
}

func New(natsCfg *config.NatsConfig, log *slog.Logger) (*Consumer, error) {
	natsURL := fmt.Sprintf("nats://%s:%s@%s:%d", natsCfg.User, natsCfg.Password, natsCfg.Host, natsCfg.Port)
	conn, err := nats.Connect(natsURL)

	if err != nil {
		log.Error("Failed to create nats consumer", err)
		panic(err)
	}

	log.Info("Created NATS consumer on", slog.String("Port", strconv.Itoa(natsCfg.Port)))

	return &Consumer{
		conn:     conn,
		log:      log,
		handlers: make(map[string]Handler),
	}, nil
}

func (c *Consumer) RegisterHandler(topic string, handler Handler) {
	c.handlers[topic] = handler
}
func (c *Consumer) Start(ctx context.Context) {
	const op = "consumer.Start"
	log := c.log.With(slog.String("op", op))

	for topic, handler := range c.handlers {
		log.Info("Subscribing to topic", slog.String("topic", topic))

		_, err := c.conn.Subscribe(topic, func(msg *nats.Msg) {
			if err := handler.Process(ctx, msg); err != nil {
				log.Error("Error processing message", slog.String("topic", topic), slog.String("err", err.Error()))
			}
		})

		if err != nil {
			log.Error("Failed to subscribe to topic", slog.String("topic", topic), slog.String("err", err.Error()))
			continue
		}
	}
	<-ctx.Done()
}

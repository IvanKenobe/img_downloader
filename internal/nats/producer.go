package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"img_downloader/internal/config"
	"log/slog"
)

type Producer struct {
	log   *slog.Logger
	conn  *nats.Conn
	topic string
}

func New(natsCfg *config.NatsConfig, topic string, log *slog.Logger) (*Producer, error) {
	natsURL := fmt.Sprintf("nats://%s:%s@%s:%d", natsCfg.User, natsCfg.Password, natsCfg.Host, natsCfg.Port)
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("nats error: %w", err)
	}
	return &Producer{conn: conn, topic: topic, log: log}, nil
}

func (p *Producer) Publish(data []byte) error {
	return p.conn.Publish(p.topic, data)
}

func (p *Producer) Close() {
	p.conn.Close()
}

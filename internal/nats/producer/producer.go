package producer

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"img_downloader/internal/config"
	"log/slog"
	"strconv"
)

type Producer struct {
	log   *slog.Logger
	conn  *nats.Conn
	topic string
}

func New(natsCfg *config.NatsConfig, topic string, log *slog.Logger) *Producer {
	natsURL := fmt.Sprintf("nats://%s:%s@%s:%d", natsCfg.User, natsCfg.Password, natsCfg.Host, natsCfg.Port)
	conn, err := nats.Connect(natsURL)

	if err != nil {
		log.Error("Failed to create nats producer", err)
		panic(err)
	}

	log.Info("Created NATS producer on", slog.String("Port", strconv.Itoa(natsCfg.Port)))
	return &Producer{conn: conn, topic: topic, log: log}
}

func (p *Producer) Publish(data []byte) error {
	return p.conn.Publish(p.topic, data)
}

func (p *Producer) Close() {
	p.conn.Close()
}

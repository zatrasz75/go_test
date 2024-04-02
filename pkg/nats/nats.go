package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
	"zatrasz75/go_test/pkg/logger"
)

type Nats struct {
	nc *nats.Conn

	allowReconnect bool
	maxReconnect   int
	reconnectWait  time.Duration
	timeout        time.Duration
	connAttempts   int
	connTimeout    time.Duration
}

func New(natsURL string, l logger.LoggersInterface, opts ...Option) (*Nats, error) {
	n := &Nats{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(n)
	}

	var err error

	for n.maxReconnect > 0 {
		// Настройка параметров повторного подключения
		var natsOpts []nats.Option

		natsOpts = append(natsOpts, nats.ReconnectWait(n.reconnectWait), nats.MaxReconnects(n.maxReconnect))

		// Добавляем параметры таймаута
		natsOpts = append(natsOpts, nats.Timeout(n.timeout))

		n.nc, err = nats.Connect(natsURL, natsOpts...)
		if err == nil {
			break
		}
		l.Info("NATS пытается подключиться, попыток осталось: %d", n.maxReconnect)

		time.Sleep(n.timeout)

		n.maxReconnect--
	}
	if err != nil {
		return nil, fmt.Errorf("nats - New - maxReconnect == 0: %w", err)
	}

	return n, nil
}

// SendLog отправляет сообщение на указанную тему subject.
func (n *Nats) SendLog(subject string, data []byte) error {
	return n.nc.Publish(subject, data)
}

// ReceiveLog получает сообщение на указанную тему subject.
func (n *Nats) ReceiveLog() (<-chan []byte, error) {
	msgChan := make(chan []byte)

	_, err := n.nc.Subscribe("logs", func(msg *nats.Msg) {
		msgChan <- msg.Data
	})
	if err != nil {
		return nil, err
	}

	return msgChan, nil
}

// Flush гарантирует, что все сообщения были обработаны сервером.
func (n *Nats) Flush() error {
	return n.nc.Flush()
}

// FlushTimeout гарантирует, что все сообщения были обработаны сервером в течение указанного времени ожидания.
func (n *Nats) FlushTimeout(timeout time.Duration) error {
	return n.nc.FlushTimeout(timeout)
}

// SubscribeToLogs подписывается на тему 'logs' и обрабатывает полученные сообщения.
func (n *Nats) SubscribeToLogs() {
	// Создание подписки на тему 'logs'
	sub, err := n.nc.SubscribeSync("logs")
	if err != nil {
		log.Fatal(err)
	}

	// Бесконечный цикл для обработки сообщений
	for {
		msg, err := sub.NextMsg(nats.DefaultTimeout)
		if err != nil {
			log.Fatal(err)
		}

		// Обработка сообщения
		n.handleLogMessage(msg)
	}
}

// handleLogMessage обрабатывает полученное сообщение.
func (n *Nats) handleLogMessage(msg *nats.Msg) {
	fmt.Printf("Received a message: %s\n", string(msg.Data))
}

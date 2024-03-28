package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Nats struct {
	nc *nats.Conn

	allowReconnect bool
	maxReconnect   int
	reconnectWait  time.Duration
	timeout        time.Duration
}

func New(natsURL string, opts ...Option) (*Nats, error) {
	n := &Nats{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(n)
	}

	var err error
	n.nc, err = nats.Connect(natsURL, nats.ReconnectWait(n.reconnectWait), nats.MaxReconnects(n.maxReconnect), nats.Timeout(n.timeout))
	if err != nil {
		return nil, err
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
			fmt.Println("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
			log.Fatal(err)
		}

		// Обработка сообщения
		n.handleLogMessage(msg)
	}
}

// handleLogMessage обрабатывает полученное сообщение.
func (n *Nats) handleLogMessage(msg *nats.Msg) {
	// Здесь вы можете обработать сообщение, например, записать его в лог или в базу данных.
	// В данном примере мы просто выводим сообщение в консоль.
	fmt.Printf("Received a message: %s\n", string(msg.Data))
}

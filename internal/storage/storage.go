package storage

import (
	"time"
	"zatrasz75/go_test/models"
)

type RepositoryInterface interface {
	// GetList Получаем всех записей с limit и offset
	GetList(limit int, offset int) ([]models.Goods, error)
	// PostList Добавляет запись с name и project_id
	PostList(g models.Goods) (models.Goods, error)
	// PatchList Обновляет запись по id и project_id
	PatchList(g models.Goods) (models.Goods, bool, error)
	// DeleteList Удаляет запись по id и project_id
	DeleteList(g models.Goods) error
	// PatchReprioritiize Обновляет priority по id и project_id у текущей записи и всех кто после +1
	PatchReprioritiize(g models.Goods) ([]models.Goods, error)
}

type RedisInterface interface {
	// GetList Получаем всех сохраненные записей
	GetList(key string) ([]models.Goods, error)
	// AddList Сохраняет все записи на определенное время
	AddList(key string, list []models.Goods) error
}

type NatsInterface interface {
	// SendLog отправляет сообщение на указанную тему subject.
	SendLog(subject string, data []byte) error
	// Flush гарантирует, что все сообщения были обработаны сервером.
	Flush() error
	// FlushTimeout гарантирует, что все сообщения были обработаны сервером в течение указанного времени ожидания.
	FlushTimeout(timeout time.Duration) error
	// SubscribeToLogs подписывается на тему 'logs' и обрабатывает полученные сообщения.
	SubscribeToLogs()
}

type ClickhouseInterface interface {
}

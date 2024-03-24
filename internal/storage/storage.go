package storage

import (
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
	Publish(subject string, data []byte) error
}

type ClickhouseInterface interface {
}

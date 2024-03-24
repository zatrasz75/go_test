package repository

import (
	"context"
	"errors"
	"zatrasz75/go_test/models"
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
	l logger.LoggersInterface
}

func New(pg *postgres.Postgres, l logger.LoggersInterface) *Store {
	return &Store{pg, l}
}

// GetList Получаем всех записей с limit и offset
func (s *Store) GetList(limit int, offset int) ([]models.Goods, error) {
	var goods []models.Goods

	// SQL запрос с LIMIT и OFFSET
	query := `SELECT id, project_id, name, description, priority, removed, created_at FROM goods ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := s.Pool.Query(context.Background(), query, limit, offset-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Goods
		err = rows.Scan(&g.ID, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		goods = append(goods, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return goods, nil
}

// PostList Добавляет запись с name и project_id
func (s *Store) PostList(g models.Goods) (models.Goods, error) {
	query := "INSERT INTO goods (project_id,name) VALUES ($1, $2) RETURNING id,project_id,name,description,priority,removed,created_at"
	row := s.Pool.QueryRow(context.Background(), query, g.ProjectId, g.Name)

	err := row.Scan(&g.ID, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return models.Goods{}, err
	}

	return g, nil
}

// PatchList Обновляет запись по id и project_id
func (s *Store) PatchList(g models.Goods) (models.Goods, bool, error) {
	// Начинаем транзакцию
	tx, err := s.Pool.Begin(context.Background())
	if err != nil {
		return models.Goods{}, false, err
	}
	defer tx.Rollback(context.Background())

	// Проверяем существование записи и блокируем её для обновления
	var exists bool
	err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE)", g.ID, g.ProjectId).Scan(&exists)
	if err != nil {
		return models.Goods{}, false, err
	}
	if !exists {
		return models.Goods{}, false, nil
	}

	// Валидируем поля перед обновлением
	if g.Name == "" {
		return models.Goods{}, false, errors.New("пустое name")
	}

	// Обновляем запись и возвращаем обновленные данные
	var updatedGoods models.Goods
	err = tx.QueryRow(
		context.Background(),
		"UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4 RETURNING id, project_id, name, description, priority, removed, created_at",
		g.Name,
		g.Description,
		g.ID,
		g.ProjectId,
	).Scan(
		&updatedGoods.ID,
		&updatedGoods.ProjectId,
		&updatedGoods.Name,
		&updatedGoods.Description,
		&updatedGoods.Priority,
		&updatedGoods.Removed,
		&updatedGoods.CreatedAt,
	)
	if err != nil {
		return models.Goods{}, false, err
	}

	// Подтверждаем транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		return models.Goods{}, false, err
	}

	return updatedGoods, true, nil
}

// DeleteList Удаляет запись по id и project_id
func (s *Store) DeleteList(g models.Goods) error {
	delet := "DELETE FROM goods WHERE id = $1 AND project_id = $2"
	_, err := s.Pool.Exec(context.Background(), delet, g.ID, g.ProjectId)
	if err != nil {
		return err
	}

	return nil
}

// PatchReprioritiize Обновляет priority по id и project_id у текущей записи и всех кто после +1
func (s *Store) PatchReprioritiize(g models.Goods) ([]models.Goods, error) {
	// Начинаем транзакцию
	tx, err := s.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	// Обновляем priority текущей записи и возвращаем измененные данные
	query := "UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3 RETURNING id, priority"
	rows, err := tx.Query(context.Background(), query, g.Priority, g.ID, g.ProjectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedGoods []models.Goods
	for rows.Next() {
		var good models.Goods
		if err := rows.Scan(&good.ID, &good.Priority); err != nil {
			return nil, err
		}
		updatedGoods = append(updatedGoods, good)
	}

	// Увеличиваем priority всех записей после текущей и возвращаем измененные данные
	query = "UPDATE goods SET priority = priority + 1 WHERE project_id = $1 AND id > $2 RETURNING id, priority"
	rows, err = tx.Query(context.Background(), query, g.ProjectId, g.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var good models.Goods
		if err = rows.Scan(&good.ID, &good.Priority); err != nil {
			return nil, err
		}
		updatedGoods = append(updatedGoods, good)
	}

	// Подтверждаем транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return updatedGoods, nil
}

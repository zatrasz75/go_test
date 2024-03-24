package models

import "time"

type Click struct {
	ID          int       `json:"ID"`
	Projectid   int       `json:"Projectid"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Priority    int       `json:"Priority"`
	Removed     bool      `json:"Removed"`
	EventTime   time.Time `json:"EventTime"`
}

type Projects struct {
	ID        int       `json:"ID"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type Goods struct {
	ID          int       `json:"id"`
	ProjectId   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"newPriority"`
	Removed     bool      `json:"Removed"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type Meta struct {
	Total   int `json:"Total"`   // Сколько всего записей
	Removed int `json:"Removed"` // Сколько записей со статусом Removed=true
	Limit   int `json:"Limit"`   // Какое ограничение стоит на вывод объектов (20)
	Offset  int `json:"Offset"`  // От какой позиции выводить данные в списке
}

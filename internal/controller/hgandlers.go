package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"zatrasz75/go_test/configs"
	"zatrasz75/go_test/internal/click"
	"zatrasz75/go_test/internal/redis"
	"zatrasz75/go_test/internal/reponats"
	"zatrasz75/go_test/internal/repository"
	"zatrasz75/go_test/internal/storage"
	"zatrasz75/go_test/models"
	"zatrasz75/go_test/pkg/logger"
)

type api struct {
	Cfg  *configs.Config
	l    logger.LoggersInterface
	repo storage.RepositoryInterface
	rd   storage.RedisInterface
	nc   storage.NatsInterface
	cl   storage.ClickhouseInterface
}

func newEndpoint(r *mux.Router, cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store, rd *redis.Store, n *reponats.Store, cl *click.Store) {
	en := &api{cfg, l, repo, rd, n, cl}

	r.HandleFunc("/", en.home).Methods(http.MethodGet)
	r.HandleFunc("/goods/list", en.goodsList).Methods(http.MethodGet)
	r.HandleFunc("/good/create", en.goodCreate).Methods(http.MethodPost)
	r.HandleFunc("/good/update", en.goodUpdate).Methods(http.MethodPatch)
	r.HandleFunc("/good/remove", en.goodRemove).Methods(http.MethodDelete)
	r.HandleFunc("/good/reprioritiize", en.goodReprioritiize).Methods(http.MethodPatch)
}

type errStruct struct{}

type httpError struct {
	Code    int
	Message string
	Details errStruct
}

func (e *httpError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, details: %v", e.Code, e.Message, e.Details)
}

func (a *api) home(w http.ResponseWriter, _ *http.Request) {
	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "https://frontend.host")

	// Выводим дополнительную строку на страницу
	str := []byte("Добро пожаловать! ")

	_, err := fmt.Fprintf(w, "<p>%s</p>", str)
	if err != nil {
		http.Error(w, "Ошибка записи на страницу", http.StatusInternalServerError)
		a.l.Error("Ошибка записи на страницу", err)
	}

	// Отправляем сообщение через NATS
	subject := "welcome"
	message := []byte("Пользователь посетил главную страницу")
	err = a.nc.Publish(subject, message)
	if err != nil {
		a.l.Error("Ошибка отправки сообщения через NATS", err)
	}
}

func (a *api) goodsList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	offsetStr := queryParams.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 1
	}
	key := fmt.Sprintf("goodsList:%d:%d", limit, offset)

	list, err := a.rd.GetList(key)
	if err != nil {
		a.l.Error("не удалось получить запись из Redis", err)
		http.Error(w, "не удалось получить запись из Redis", http.StatusInternalServerError)
		return
	}

	if list == nil {
		list, err = a.repo.GetList(limit, offset)
		if err != nil {
			a.l.Error("не получилось получить данные из goods", err)
			http.Error(w, "не получилось получить данные из goods", http.StatusInternalServerError)
			return
		}

		err = a.rd.AddList(key, list)
		if err != nil {
			a.l.Error("не удалось сделать запись в Redis", err)
			http.Error(w, "не удалось сделать запись в Redis", http.StatusInternalServerError)
			return
		}

		for i, v := range list {
			k := strconv.Itoa(v.ID)
			singleItemList := []models.Goods{list[i]}
			err = a.rd.AddList(k, singleItemList)
			if err != nil {
				a.l.Error("не удалось сделать запись в Redis", err)
				http.Error(w, "не удалось сделать запись в Redis", http.StatusInternalServerError)
				return
			}
		}
	}

	var removed int
	for _, v := range list {
		if v.Removed == true {
			removed += 1
		}
	}

	data := struct {
		Meta  models.Meta
		Goods []models.Goods
	}{
		Meta: models.Meta{
			Total:   len(list),
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: list,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "ошибка при форматировании данных", http.StatusInternalServerError)
		a.l.Error("ошибка при форматировании данных: ", err)
		return
	}

	// Устанавливаем правильный Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

func (a *api) goodCreate(w http.ResponseWriter, r *http.Request) {
	var g models.Goods

	queryParams := r.URL.Query()
	idStr := queryParams.Get("projectId")
	projectId, err := strconv.Atoi(idStr)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	g.ProjectId = projectId

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&g); err != nil {
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		a.l.Error("не удалось проанализировать запрос JSON", err)

	}

	list, err := a.repo.PostList(g)
	if err != nil {
		a.l.Error("Ошибка при добавления данных", err)
		http.Error(w, "Ошибка при добавления данных", http.StatusInternalServerError)
		return
	}

	jsonList, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "ошибка при форматировании данных", http.StatusInternalServerError)
		a.l.Error("ошибка при форматировании данных: ", err)
		return
	}

	// Устанавливаем правильный Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonList)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

func (a *api) goodUpdate(w http.ResponseWriter, r *http.Request) {
	var g models.Goods

	queryParams := r.URL.Query()
	project := queryParams.Get("projectId")
	idStr := queryParams.Get("id")
	projectId, err := strconv.Atoi(project)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	g.ID = id
	g.ProjectId = projectId

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&g); err != nil {
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		a.l.Error("не удалось проанализировать запрос JSON", err)
		return
	}
	if g.Name == "" {
		http.Error(w, "не удалось получить имя", http.StatusNotFound)
		a.l.Debug("не удалось получить имя")
		return
	}

	list, exists, err := a.repo.PatchList(g)
	if err != nil {
		http.Error(w, "не удалось обновить запись", http.StatusBadRequest)
		a.l.Error("не удалось обновить запись", err)
		return
	}
	if !exists {
		http.Error(w, "нет такой записи", http.StatusNotFound)
		a.l.Error("нет такой записи", err)
		httpErr := &httpError{
			Code:    3,
			Message: "errors.good.notFound",
			Details: errStruct{},
		}
		// Отправляем ошибку клиенту
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	} else {
		k := strconv.Itoa(list.ID)
		singleItemList := []models.Goods{list}
		err = a.rd.AddList(k, singleItemList)
		if err != nil {
			a.l.Error("не удалось сделать запись в Redis", err)
			http.Error(w, "не удалось сделать запись в Redis", http.StatusInternalServerError)
			return
		}
	}

	jsonList, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "ошибка при форматировании данных", http.StatusInternalServerError)
		a.l.Error("ошибка при форматировании данных: ", err)
		return
	}

	// Устанавливаем правильный Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonList)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

func (a *api) goodRemove(w http.ResponseWriter, r *http.Request) {
	var g models.Goods

	queryParams := r.URL.Query()
	project := queryParams.Get("projectId")
	idStr := queryParams.Get("id")
	if project == "" || idStr == "" {
		http.Error(w, "не удалось получить запись", http.StatusNotFound)
		a.l.Debug("не удалось получить запись")
		return
	}
	projectId, err := strconv.Atoi(project)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	g.ID = id
	g.ProjectId = projectId

	err = a.repo.DeleteList(g)
	if err != nil {
		if err != nil {
			a.l.Error("не получилось удалить данные из goods", err)
			http.Error(w, "не получилось удалить данные из goods", http.StatusInternalServerError)
			return
		}
	}
	data := struct {
		Id         int
		CampaignId int
		Removed    bool
	}{
		Id:         id,
		CampaignId: id,
		Removed:    true,
	}
	dataJason, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "ошибка при форматировании данных", http.StatusInternalServerError)
		a.l.Error("ошибка при форматировании данных: ", err)
		return
	}

	// Устанавливаем правильный Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dataJason)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

func (a *api) goodReprioritiize(w http.ResponseWriter, r *http.Request) {
	var g models.Goods

	queryParams := r.URL.Query()
	project := queryParams.Get("projectId")
	idStr := queryParams.Get("id")
	if project == "" || idStr == "" {
		http.Error(w, "не удалось получить запись", http.StatusNotFound)
		a.l.Debug("не удалось получить запись")
		return
	}
	projectId, err := strconv.Atoi(project)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
	}
	g.ID = id
	g.ProjectId = projectId

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&g); err != nil {
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		a.l.Error("не удалось проанализировать запрос JSON", err)
		return
	}

	list, err := a.repo.PatchReprioritiize(g)
	if err != nil {
		if err != nil {
			http.Error(w, "не удалось обновить запись", http.StatusBadRequest)
			a.l.Error("не удалось обновить запись", err)
			return
		}
	}

	response := struct {
		Priorities []struct {
			ID       int `json:"id"`
			Priority int `json:"priority"`
		} `json:"priorities"`
	}{
		Priorities: make([]struct {
			ID       int `json:"id"`
			Priority int `json:"priority"`
		}, len(list)),
	}
	for i, good := range list {
		response.Priorities[i] = struct {
			ID       int `json:"id"`
			Priority int `json:"priority"`
		}{
			ID:       good.ID,
			Priority: good.Priority,
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "не удалось преобразовать ответ в JSON", http.StatusInternalServerError)
		a.l.Error("не удалось преобразовать ответ в JSON", err)
		return
	}

	// Устанавливаем правильный Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

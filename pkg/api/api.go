package api

import (
	"encoding/json"
	"net/http"
	storage "newsgo/pkg/storage"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r  *mux.Router       // маршрутизатор запросов
	db storage.Interface // база данных
}

// Конструктор API.
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов. ??? нужен ли метод для добавления в базу?
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
	// добавление новости в базу
	api.r.HandleFunc("/news", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)
}

// postsHandler возвращает все публикации.
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := api.db.Posts(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

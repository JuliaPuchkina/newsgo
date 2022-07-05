package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	api "newsgo/pkg/api"
	"newsgo/pkg/rss"
	storage "newsgo/pkg/storage"
	db "newsgo/pkg/storage/db"
	"os"
	"time"
)

// сервер newsgo
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// объект сервера
	var srv server

	// объект базы данных postgresql
	db, err := db.New("postgres://postgres:postgres@127.0.0.1/newsag")
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// каналы для обработки новостей и ошибок
	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)

	// получение и парсинг ссылок из конфига
	myLinks := getRss("config.json", chanErrs)
	for i := range myLinks.LinkArr {
		go parseNews(myLinks.LinkArr[i], db, chanErrs, chanPosts)
	}

	// обработка потока новостей
	go func() {
		for posts := range chanPosts {
			for i := range posts {
				db.AddPost(posts[i])
			}
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range chanErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

type Links struct {
	LinkArr []string `json:"rss"`
}

// функция, получающая отдельные ссылки из конфига, ошибки направляются в поток ошибок
func getRss(fileName string, errors chan<- error) Links {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		errors <- err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var links Links

	json.Unmarshal(byteValue, &links)

	return links
}

// получение новостей по ссылкам и отправка новостей и ошибок в соответствующие каналы
func parseNews(link string, db *db.Store, errs chan<- error, posts chan<- []storage.Post) {
	for {
		newsPosts, err := rss.RssToStruct(link)
		if err != nil {
			errs <- err
			continue
		}
		posts <- newsPosts
		time.Sleep(time.Minute * 5)
	}
}

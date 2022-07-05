package db

import (
	storage "newsgo/pkg/storage"
	"testing"
)

func TestStore_AddPost(t *testing.T) {
	dataBase, err := New("postgres://postgres:postgres@127.0.0.1/newsag")
	post := storage.Post{
		Title:   "Unit Test Task",
		Content: "Task Content",
		PubTime: 5,
		Link:    "link",
	}
	id := dataBase.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Создана задача с id:", id)
}

// этот тест долго думала, как писать, и в итоге, мне кажется, получилось что-то не то
func TestStore_Posts(t *testing.T) {
	dataBase, _ := New("postgres://postgres:postgres@127.0.0.1/newsag")
	news, err := dataBase.Posts(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)

}

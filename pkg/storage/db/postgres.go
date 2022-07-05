package db

import (
	"context"
	storage "newsgo/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// конструктор объекта хранилища
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// Posts выводит все существующие публикации
func (s *Store) Posts(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := s.db.Query(context.Background(), `
	SELECT id, title, content, published, link FROM news
	ORDER BY published DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// AddPost создает новую публикацию
func (s *Store) AddPost(p storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO news (title, content, published, link)
		VALUES ($1, $2, $3, $4);
		`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	).Scan()
	return err
}

DROP TABLE IF EXISTS news;

-- новости
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
	title TEXT, -- заголовок новости
	content TEXT NOT NULL UNIQUE, -- текст 
    published BIGINT DEFAULT 0, -- время публикации новости
    link TEXT -- ссылка на новость    
);
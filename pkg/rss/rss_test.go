package rss

import "testing"

func TestRssToStruct(t *testing.T) {
	posts, err := RssToStruct("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("данные не рскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(posts), posts)
}

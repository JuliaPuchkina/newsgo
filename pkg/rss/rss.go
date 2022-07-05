package rss

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	storage "newsgo/pkg/storage"
	"strings"
	"time"
)

// структура данных, получаемых из rss
type MyXMLstruct struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title   string `xml:"channel>title"`
	Link    string `xml:"channel>link"`
	Content string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []Item `xml:"channel>item"`
}

// структура для отдельного поста
type Item struct {
	// Required
	Title   string        `xml:"title"`
	Link    string        `xml:"link"`
	Content template.HTML `xml:"description"`
	PubDate string        `xml:"pubDate"`
}

// преобразование полученных xml данных в заданную структуру, затем в массив новостей
func RssToStruct(link string) ([]storage.Post, error) {
	var posts MyXMLstruct
	if xmlBytes, err := getXML(link); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &posts)

	}

	var news []storage.Post
	for j := range posts.ItemList {
		var item storage.Post
		item.Title = posts.ItemList[j].Title
		item.Content = string(posts.ItemList[j].Content)
		item.Link = posts.ItemList[j].Link

		posts.ItemList[j].PubDate = strings.ReplaceAll(posts.ItemList[j].PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", posts.ItemList[j].PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", posts.ItemList[j].PubDate)
		}
		if err == nil {
			item.PubTime = t.Unix()
		}
		news = append(news, item)
	}

	return news, nil
}

// получение xml данных по ссылке
func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

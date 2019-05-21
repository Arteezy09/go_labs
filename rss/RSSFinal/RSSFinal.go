package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"sort"

	"github.com/RealJK/rss-parser-go"
	//"strings"  // пакет для работы с  UTF-8 строками
)

type byDate []rss.Item

func (s byDate) Len() int {
	return len(s)
}
func (s byDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDate) Less(i, j int) bool {
	return s[i].PubDate < s[j].PubDate
}

func outputRSSItems(w http.ResponseWriter, items []rss.Item) {
	sort.Sort(byDate(items))
	fmt.Fprintf(w, "<html><head></head><body>")
	for v := range items {
		item := items[v]
		fmt.Fprintf(w, "<div style='margin-bottom:20px;'>")
		//fmt.Fprintf(w, "%s\n", item.Source)
		fmt.Fprintf(w, "%s<br>", item.PubDate)
		fmt.Fprintf(w, "%s\n", item.Title)
		fmt.Fprintf(w, "<details><summary>Подробнее</summary><p>%s</p></details>", item.Description)
		fmt.Fprintf(w, "</div>")
	}
	fmt.Fprintf(w, "</body></html>")
}
func outputRSSByURL(w http.ResponseWriter, url string) {
	rssObject, error := rss.ParseRSS(url)
	if error != nil {
		outputRSSItems(w, rssObject.Channel.Items)
	}
}
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<a href=/>Всё</a> ")
	fmt.Fprintf(w, "<a href=/rss0>Интернет</a> ")
	fmt.Fprintf(w, "<a href=/rss1>Экономика</a> ")
	fmt.Fprintf(w, "<a href=/rss2>Теннис</a> ")
	r.ParseForm() //анализ аргументов
	urls := [3]string{"https://news.yandex.ru/internet.rss",
		"https://news.yandex.ru/business.rss",
		"https://news.yandex.ru/tennis.rss"}
	if r.URL.Path == "/" {
		allItems := []rss.Item{}

		for urlI := range urls {
			url := urls[urlI]
			parsedRSS, error := rss.ParseRSS(url)
			if error != nil {
				for itemI := range parsedRSS.Channel.Items {
					item := parsedRSS.Channel.Items[itemI]
					allItems = append(allItems, item)
				}
			}
		}
		outputRSSItems(w, allItems)
	} else if r.URL.Path == "/rss0" {
		outputRSSByURL(w, urls[0])
	} else if r.URL.Path == "/rss1" {
		outputRSSByURL(w, urls[1])
	} else if r.URL.Path == "/rss2" {
		outputRSSByURL(w, urls[2])
	}
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":9008", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

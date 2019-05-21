package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола

	"github.com/RealJK/rss-parser-go"
	//"strings"  // пакет для работы с  UTF-8 строками
)

func outputRSSByURL(w http.ResponseWriter, url string) {
	rssObject, error := rss.ParseRSS(url)

	if error != nil {
		fmt.Fprintf(w, "<html><head></head><body>")
		for v := range rssObject.Channel.Items {
			item := rssObject.Channel.Items[v]
			fmt.Fprintf(w, "<div style='margin-bottom:20px;'>")
			fmt.Fprintf(w, "%s\n", item.Title)
			fmt.Fprintf(w, "<details><summary>Подробнее</summary><p>%s</p></details>", item.Description)
			fmt.Fprintf(w, "</div>")
		}
		fmt.Fprintf(w, "</body></html>")
	}
}
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов

	if r.URL.Path == "/" {
		rssObject0, _ := rss.ParseRSS("http://blagnews.ru/rss_vk.xml")
		rssObject1, _ := rss.ParseRSS("https://news.yandex.ru/science.rss")
		rssObject2, _ := rss.ParseRSS("https://news.yandex.ru/realty.rss")

		allItems := []rss.Item{}
		for v := range rssObject0.Channel.Items {
			allItems = append(allItems, rssObject0.Channel.Items[v])
		}
		for v := range rssObject1.Channel.Items {
			allItems = append(allItems, rssObject1.Channel.Items[v])
		}
		for v := range rssObject2.Channel.Items {
			allItems = append(allItems, rssObject2.Channel.Items[v])
		}

		fmt.Fprintf(w, "<html><head></head><body>")
		for v := range allItems {
			item := allItems[v]
			fmt.Fprintf(w, "<div style='margin-bottom:20px;'>")
			fmt.Fprintf(w, "%s\n", item.Title)
			fmt.Fprintf(w, "<details><summary>Подробнее</summary><p>%s</p></details>", item.Description)
			fmt.Fprintf(w, "</div>")
		}
		fmt.Fprintf(w, "</body></html>")
	} else if r.URL.Path == "/rss0" {
		outputRSSByURL(w, "http://blagnews.ru/rss_vk.xml")
	} else if r.URL.Path == "/rss1" {
		outputRSSByURL(w, "https://news.yandex.ru/science.rss")
	} else if r.URL.Path == "/rss2" {
		outputRSSByURL(w, "https://news.yandex.ru/realty.rss")
	}
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":9008", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

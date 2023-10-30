package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kouhei-github/book-advertisement-site/crontab/cron_crawler"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type ZenRequest struct {
	Page int `json:"page"`
}

func BookInZennHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"status": "Mthod Not Allowed", "Text": "失敗"})
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "Internal Server Error", "Text": "失敗"})
		return
	}
	var requestBody ZenRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "Internal Server Error", "Text": "Faild json to struct"})
		return
	}

	zenn := "https://zenn.dev/topics/%E6%8A%80%E8%A1%93%E6%9B%B8"
	var result []string
	var wg sync.WaitGroup
	re := regexp.MustCompile(`<a class="ArticleList_link__4Igs4" href="(.+?)">`)

	for i := 1; i < 100; i++ {
		wg.Add(1)
		url := zenn + "?page=" + strconv.Itoa(i)
		article := cron_crawler.ZennArticle{Url: url}
		go func(url string, imp cron_crawler.BookSearchImp) {
			defer wg.Done()
			html, err := imp.Run()

			if err != nil {
				fmt.Println(err.Error())
				return
			}
			matches := re.FindAllStringSubmatch(html, -1)

			if len(matches) == 0 {
				return
			}
			for _, match := range matches {
				detailUrl := "https://zenn.dev" + match[1]
				books, err := imp.FindAmazonBook(detailUrl)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				result = append(result, books...)
			}

		}(url, article)

	}
	wg.Wait()

	//スプレッドに書き込み
	json.NewEncoder(w).Encode(map[string][]string{"result": result})
}

func BookInQiitaHandler(w http.ResponseWriter, r *http.Request) {
	var includeBookUrl []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 1; i < 100; i++ {
		wg.Add(1)
		page := strconv.Itoa(i)
		url := "http://qiita.com/api/v2/items?query=tag:技術書&page=" + page + "&per_page=100"
		qiita := cron_crawler.Qiita{Url: url}
		//go func(qiita cron_crawler.Qiita) {
		//	defer wg.Done()
		//	articles, err := qiita.Run()
		//	if err != nil {
		//		return
		//	}
		//	fmt.Println(len(articles))
		//	if len(articles) != 0 {
		//		mu.Lock()
		//		for _, article := range articles {
		//			includeBookUrl = append(includeBookUrl, article.URL)
		//		}
		//		mu.Unlock()
		//	}
		//
		//}(qiita)
		defer wg.Done()
		articles, _ := qiita.Run()

		fmt.Println(len(articles))
		if len(articles) != 0 {
			mu.Lock()
			for _, article := range articles {
				includeBookUrl = append(includeBookUrl, article.URL)
			}
			mu.Unlock()
		}

	}
	wg.Wait()
	fmt.Println("len(includeBookUrl)")

	//fmt.Println(len(includeBookUrl))
}

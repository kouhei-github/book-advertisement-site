package cron_crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Qiita struct {
	Url string
}

type QiitaArticle struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Tags      []Tag     `json:"tags"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}
type Tag struct {
	Name string `json:"name"`
}

func (e Qiita) Run() ([]QiitaArticle, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", e.Url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var qiitaArticle []QiitaArticle
	if err = json.Unmarshal(body, &qiitaArticle); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return qiitaArticle, nil
}

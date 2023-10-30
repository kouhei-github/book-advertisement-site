package cron_crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type ZennArticle struct {
	Url string
}

var re = regexp.MustCompile(`<a href="(.+?)"`)

func (e ZennArticle) Run() (string, error) {
	return getHttpRequest(e.Url)
}

func (e ZennArticle) FindAmazonBook(url string) ([]string, error) {
	html, err := getHttpRequest(url)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(html, -1)
	var result []string
	if len(matches) != 0 {
		for _, match := range matches {
			if !strings.Contains(match[1], "www.amazon.co.jp") {
				continue
			}
			result = append(result, match[1]) // グループ1（"(...)"部分）を抽出
		}
	}

	return result, nil
}

func getHttpRequest(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

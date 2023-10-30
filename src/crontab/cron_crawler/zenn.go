package cron_crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
)

type Zenn struct {
	Url string
}

type ZennArticlesResponse struct {
	Articles []zennArticle `json:"articles"`
}

type zennArticle struct {
	Slug string `json:"slug"`
}

// Run はZennから最新の記事を取得し、それらをリポジトリ層で保存します。
// まずHTTPクライアントを作成し、指定したURLに対してGETリクエストを行います。
// エラーが発生した場合はエラーメッセージを表示し、関数を終了します。
// 次にレスポンスボディを読み取り、それをZennArticlesResponse構造体にアンマーシャルします。
// エラーが発生した場合もエラーメッセージを表示し、関数を終了します。
//
// 次に、取得した各記事について、それぞれの詳細ページを取得します。
// この際、記事に関連付けられた各タグも生成されます。
// タグはリポジトリ層で検索し、存在しない場合は新たに保存します。
// エラーが発生した場合はエラーメッセージを表示します。
//
// タグの生成が完了したら、新たにArticle構造体を作成し、記事情報をセットします。
// これをarticlesスライスに追加します。
//
// 最後に、取得した各記事をリポジトリ層で保存します。
// ここでもエラーメッセージを表示し、関数を終了します。
func (e Zenn) Run() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", e.Url+"?order=latest", nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var zenArticlesResponse ZennArticlesResponse
	if err := json.Unmarshal(body, &zenArticlesResponse); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(zenArticlesResponse)
}

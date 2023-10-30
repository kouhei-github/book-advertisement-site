package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/kouhei-github/book-advertisement-site/crontab"
	"github.com/kouhei-github/book-advertisement-site/route"
	"github.com/kouhei-github/book-advertisement-site/utils/cors"

	"net/http"
	"os"
)

func main() {
	// crontabでジョブの実行
	crontab.ToStartCron()

	// ルーターの設定
	router := route.Router{Mutex: http.NewServeMux()}
	router.GetRouter()

	// CORS (Cross Origin Resource Sharing)の設定
	// アクセスを許可するドメイン等を設定します
	corsOrigin := cors.NewCorOrigin()

	handler := corsOrigin.Handler(router.Mutex)
	// Webサーバー起動時のエラーハンドリング => localhostの時コメントイン必要
	if os.Getenv("ENVIRONMENT") == "local" {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			panic(err)
		}
	}

	// AWS Lambdaとの連携設定
	lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
}

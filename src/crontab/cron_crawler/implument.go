package cron_crawler

type BookSearchImp interface {
	Run() (string, error)
	FindAmazonBook(url string) ([]string, error)
}

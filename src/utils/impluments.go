package utils

type SpreadReader interface {
	Get(sheetRange string) ([][]interface{}, error)
}

type SpreadWriter interface {
	Write(sheetRange string, updateValues [][]interface{}) error
}

package utils

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"os"
)

type SpreadSheet struct {
	Srv           *sheets.Service
	SpreadSheetId string
	Context       context.Context
}

func getServiceAccountJson() map[string]string {
	return map[string]string{
		"type":                        os.Getenv("TYPE"),
		"project_id":                  os.Getenv("PROJECT_ID"),
		"private_key_id":              os.Getenv("PROJECT_KEY_ID"),
		"private_key":                 os.Getenv("PRIVATE_KEY"),
		"client_email":                os.Getenv("CLIENT_EMAIL"),
		"client_id":                   os.Getenv("CLIENT_ID"),
		"auth_uri":                    os.Getenv("AUTH_URI"),
		"token_uri":                   os.Getenv("TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("AUTH_PROVIDER_X506_CERT_URI"),
		"client_x509_cert_url":        os.Getenv("CLIENT_X506_CERT_URI"),
	}
}

func NewSpreadSheetReader(spreadId string) (SpreadReader, error) {
	serviceAccount := getServiceAccountJson()
	bytes, err := json.Marshal(serviceAccount)
	if err != nil {
		return SpreadSheet{}, err
	}

	config, err := google.JWTConfigFromJSON(bytes, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return SpreadSheet{}, err
	}
	context.Background()
	ctx := context.Background()
	client := config.Client(ctx)

	sheetService, err := sheets.New(client)
	if err != nil {
		return SpreadSheet{}, err
	}

	return SpreadSheet{Srv: sheetService, Context: ctx, SpreadSheetId: spreadId}, nil
}

func NewSpreadSheetWriter(spreadId string) (SpreadWriter, error) {
	serviceAccount := getServiceAccountJson()

	bytes, err := json.Marshal(serviceAccount)
	if err != nil {
		return SpreadSheet{}, err
	}

	config, err := google.JWTConfigFromJSON(bytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return SpreadSheet{}, err
	}
	context.Background()
	ctx := context.Background()
	client := config.Client(ctx)

	sheetService, err := sheets.New(client)
	if err != nil {
		return SpreadSheet{}, err
	}

	return SpreadSheet{Srv: sheetService, Context: ctx, SpreadSheetId: spreadId}, nil
}

func (spread SpreadSheet) Get(sheetRange string) ([][]interface{}, error) {
	spreadData, err := spread.Srv.Spreadsheets.Values.Get(spread.SpreadSheetId, sheetRange).Do()
	return spreadData.Values, err
}

func (spread SpreadSheet) Write(sheetRange string, updateValues [][]interface{}) error {
	updateRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         updateValues,
	}
	_, err := spread.Srv.Spreadsheets.Values.Update(
		spread.SpreadSheetId, sheetRange, updateRange).ValueInputOption(
		"USER_ENTERED").Context(spread.Context).Do()
	return err
}

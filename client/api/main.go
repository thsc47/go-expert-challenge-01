package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuotationResponse struct {
	Bid string `json:"bid"`
}

func main() {
	quotation, err := getQuotation()
	if err != nil {
		panic(err)
	}

	err = SaveQuotation(quotation)
	if err != nil {
		panic(err)
	}
}

func getQuotation() (*QuotationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Erro ao montar requisição para obter a cotação: %s", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("Erro na comunicação com servidor, tempo de resposta excedido")
			return nil, err
		}

		fmt.Printf("Erro ao obter a cotação: %s", err.Error())
		return nil, err
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Erro ao obter a cotação, status code: %d, body: %s\n", resp.StatusCode, string(body))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao converter resultado da cotação: %s", err.Error())
		return nil, err
	}

	var quotation QuotationResponse
	err = json.Unmarshal(body, &quotation)
	if err != nil {
		fmt.Printf("Erro ao ler resultado da cotação: %s", err.Error())
		return nil, err
	}

	return &quotation, nil
}

func closeBody(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		fmt.Printf("Erro ao fechar conexão com response.body: %s", err.Error())
		return
	}
}

func SaveQuotation(quotation *QuotationResponse) error {
	filePath := "../quotation.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "./quotation.txt"
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro abrir file quotation.txt para registro: %s", err.Error())
		return err
	}
	defer closeFile(file)

	err = os.WriteFile(filePath, []byte(fmt.Sprintf("Dólar: %s\n", quotation.Bid)), os.ModePerm)
	if err != nil {
		fmt.Printf("Erro ao registrar cotação no file quotation.txt: %s", err.Error())
		return err
	}

	fmt.Printf("Cotação do dólar registrada com sucesso no file quotation.txt")
	return nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Printf("Erro ao fechar conexão com arquivo.txt: %s", err.Error())
		return
	}
}

package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar requisição:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro ao fazer requisição:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Erro na resposta do servidor (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Fatal("Erro ao decodificar resposta:", err)
	}
	
	err = os.WriteFile("cotacao.txt", []byte("Dólar: "+cotacao.Bid+"\n"), 0644)
	if err != nil {
		log.Fatal("Erro ao salvar cotacao.txt:", err)
	}

	log.Println("Cotação salva com sucesso em cotacao.txt:", cotacao.Bid)
}

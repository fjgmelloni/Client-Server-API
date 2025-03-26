package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Bid  string `json:"bid"`
}

type CotacaoCompleta struct {
	USDBRL Cotacao `json:"USDBRL"`
}

type CotacaoDB struct {
	ID    uint      `gorm:"primaryKey"`
	Code  string
	Name  string
	Bid   string
	Data  time.Time
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancelAPI()

	cotacao, err := BuscarCotacao(ctxAPI)
	if err != nil {
		http.Error(w, "Erro ao buscar cotação: "+err.Error(), http.StatusRequestTimeout)
		log.Println("Erro ao buscar cotação:", err)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelDB()

	if err := GravarCotacao(ctxDB, cotacao); err != nil {
		log.Println("Erro ao salvar no banco:", err)
		// continua execução mesmo com erro de banco
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": cotacao.USDBRL.Bid})
}

func BuscarCotacao(ctx context.Context) (*CotacaoCompleta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao CotacaoCompleta
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func GravarCotacao(ctx context.Context, c *CotacaoCompleta) error {
	db, err := gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.WithContext(ctx).AutoMigrate(&CotacaoDB{}); err != nil {
		return err
	}

	cotacao := CotacaoDB{
		Code: c.USDBRL.Code,
		Name: c.USDBRL.Name,
		Bid:  c.USDBRL.Bid,
		Data: time.Now(),
	}

	if err := db.WithContext(ctx).Create(&cotacao).Error; err != nil {
		return err
	}

	log.Println("Cotação salva com sucesso:", cotacao.Bid)
	return nil
}

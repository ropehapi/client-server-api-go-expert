package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Error struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	defer res.Body.Close()
	respInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	response := make(map[string]map[string]interface{})
	_ = json.Unmarshal(respInBytes, &response)

	bid := response["USDBRL"]["bid"]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)

	//TODO: Continuar a partir daqui
	db, err := sql.Open("sqlite3", "client_server_api_go_expert_db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	defer db.Close()

	_, err = db.Exec("Insert into cotacao (data_hora, valor) values (1997-03-01, 5.40)")
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
}

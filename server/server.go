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

	//Em tese, isso deve estar funcionando, mas nao tenho certeza
	//pois não consegui testar devido a um erro de compilação do meu OS
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("Insert into cotacao (valor, data) values ($1, $2)")
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), bid)
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	error := Error{Message: "Deu certo"}
	json.NewEncoder(w).Encode(error)
}

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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println(err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	response := make(map[string]string)
	_ = json.Unmarshal(bytes, &response)

	bid := response["bid"]

	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err.Error())
	}

	stringGravar := "Dólar: " + string(bid)
	_, err = f.Write([]byte(stringGravar))
	if err != nil {
		log.Println(err.Error())
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Board struct {
	Id   string
	Name string
}

func main() {
	baseURL := "https://api.trello.com/1/members/me/boards"
	key := os.Getenv("KEY")
	token := os.Getenv("TOKEN")
	query := url.Values{}
	query.Set("filter", "open")
	query.Set("key", key)
	query.Set("token", token)
	res, err := http.Get(baseURL + "?" + query.Encode())
	if err != nil {
		log.Fatalf("oops1: %s", err)
	}
	defer res.Body.Close()
	boards := make([]Board, 0)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("oops2: %s", err)
		return
	}
	json.Unmarshal(bytes, &boards)
	if err != nil {
		log.Fatalf("oops3: %s", err)
		return
	}
	for _, board := range boards {
		fmt.Printf("%s (%s)\n", board.Name, board.Id)
	}
	// _, err = io.Copy(os.Stdout, res.Body)
	// if err != nil {
	// 	fmt.Printf("Error copying response body: %s", err)
	// 	return
	// }
}

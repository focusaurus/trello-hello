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

type trello struct {
	baseURL string
}

func (t *trello) ListBoards() ([]Board, error) {
	key := os.Getenv("KEY")
	token := os.Getenv("TOKEN")
	query := url.Values{}
	query.Set("filter", "open")
	query.Set("key", key)
	query.Set("token", token)
	res, err := http.Get(t.baseURL + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	boards := make([]Board, 0)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &boards)
	return boards, err
}

func newTrello(baseURL string) trello {
	return trello{baseURL}
}

func main() {
	baseURL := "https://api.trello.com/1/members/me/boards"
	trello := newTrello(baseURL)
	boards, err := trello.ListBoards()
	if err != nil {
		log.Fatalf("oops1: %s", err)
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

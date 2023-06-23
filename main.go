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
	Id    string
	Name  string
	Lists []List
}

type List struct {
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
	query.Set("fields", "id,name,lists")
	query.Set("lists", "open")
	query.Set("list_fields", "id,name")
	query.Set("key", key)
	query.Set("token", token)
	res, err := http.Get(t.baseURL + "/1/members/me/boards?" + query.Encode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListBoards err1: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	boards := make([]Board, 0)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListBoards err2: %s", err)
		return nil, err
	}
	// os.Stdout.Write(bytes)
	err = json.Unmarshal(bytes, &boards)
	return boards, err
}

func (t *trello) ListLists(board Board) ([]List, error) {
	key := os.Getenv("KEY")
	token := os.Getenv("TOKEN")
	query := url.Values{}
	query.Set("filter", "open")
	query.Set("key", key)
	query.Set("token", token)
	query.Set("cards", "open")
	query.Set("card_fields", "name")
	res, err := http.Get(t.baseURL + "/1/boards/" + board.Id + "/lists?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	lists := make([]List, 0)
	bytes, err := io.ReadAll(res.Body)
	os.Stdout.Write(bytes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &lists)
	return lists, err
}

func newTrello(baseURL string) trello {
	return trello{baseURL}
}

func main() {
	baseURL := "https://api.trello.com"
	trello := newTrello(baseURL)
	boards, err := trello.ListBoards()
	if err != nil {
		log.Fatalf("oops1: %s", err)
	}
	for _, board := range boards {
		fmt.Printf("📋%s (%s)\n", board.Name, board.Id)
		// lists, err := trello.ListLists(board)
		// if err != nil {
		// 	log.Fatalf("oops2: %s", err)
		// }
		for _, list := range board.Lists {
			fmt.Printf("\t📃%s (%s)\n", list.Name, list.Id)
		}
	}
	// _, err = io.Copy(os.Stdout, res.Body)
	// if err != nil {
	// 	fmt.Printf("Error copying response body: %s", err)
	// 	return
	// }
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

type Card struct {
	Id   string
	Name string
}

type trello struct {
	baseURL string
	key     string
	token   string
}

func (t *trello) getJSON(path string, query url.Values, decodeTo any) error {
	query.Set("key", t.key)
	query.Set("token", t.token)
	res, err := http.Get(t.baseURL + path + "?" + query.Encode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "getJSON err1: %s", err)
		return err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getJSON err2: %s", err)
		return err
	}
	// os.Stdout.Write(bytes)
	return json.Unmarshal(bytes, decodeTo)
}

func (t *trello) ListBoards() ([]Board, error) {
	query := url.Values{}
	query.Set("filter", "open")
	query.Set("fields", "id,name,lists")
	query.Set("lists", "open")
	query.Set("list_fields", "id,name")

	path := "/1/members/me/boards"
	boards := make([]Board, 0)
	err := t.getJSON(path, query, &boards)
	return boards, err
}

func (t *trello) ListCards(list List) ([]Card, error) {
	query := url.Values{}
	path := "/1/lists/" + list.Id + "/cards"
	cards := make([]Card, 0)
	err := t.getJSON(path, query, &cards)
	return cards, err
}

func newTrello(baseURL string) trello {
	return trello{
		baseURL: baseURL,
		key:     os.Getenv("KEY"),
		token:   os.Getenv("TOKEN"),
	}
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
		for _, list := range board.Lists {
			if matched, _ := regexp.MatchString("(To Do|Doing)", list.Name); !matched {
				// if strings.Contains(list.Name, "Done") {
				continue
			}
			fmt.Printf("  📃%s (%s)\n", list.Name, list.Id)
			cards, err := trello.ListCards(list)
			if err != nil {
				log.Fatalf("oops2: %s", err)
			}
			for _, card := range cards {
				fmt.Printf("    🪧%s (%s)\n", card.Name, card.Id)
			}
		}
	}
}

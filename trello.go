package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/go-playground/validator/v10"
)

type Row struct {
	ID   string `json:"id"`
	Name string
}

func (r Row) String() string {
	return r.Name
}

type Board struct {
	Row
	Lists []Row
}

type trelloClient struct {
	BaseURL string `validate:"min=7"`
	Key     string `validate:"min=20"`
	Token   string `validate:"min=50"`
}

func (t *trelloClient) getJSON(path string, query url.Values, decodeTo any) error {
	query.Set("key", t.Key)
	query.Set("token", t.Token)
	res, err := http.Get(t.BaseURL + path + "?" + query.Encode())
	if err != nil {
		return fmt.Errorf("error in getJSON http.Get %w", err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error from Trello API: %s", string(bytes))
	}
	if err != nil {
		return fmt.Errorf("error in getJSON io.ReadAll %w", err)
	}
	// Useful for debugging
	// os.Stdout.Write(bytes)
	err = json.Unmarshal(bytes, decodeTo)
	if err != nil {
		return fmt.Errorf("invalid JSON from Trello API: %s", err)
	}
	return nil
}

func (t *trelloClient) ListBoards() ([]Board, error) {
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

func (t *trelloClient) ListCards(list Row) ([]Row, error) {
	query := url.Values{}
	path := "/1/lists/" + list.ID + "/cards"
	cards := make([]Row, 0)
	err := t.getJSON(path, query, &cards)
	return cards, err
}

func newTrello(baseURL string) (*trelloClient, error) {
	t := &trelloClient{
		BaseURL: baseURL,
		Key:     os.Getenv("KEY"),
		Token:   os.Getenv("TOKEN"),
	}
	validate := validator.New()
	err := validate.Struct(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

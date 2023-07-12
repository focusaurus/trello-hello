package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/go-playground/validator/v10"
)

type Row struct {
	ID   string
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
	BaseURL string `validate:"http_url"`
	Key     string `validate:"min=20"`
	Token   string `validate:"min=50"`
}

func (t *trelloClient) getJSON(ctx context.Context, path string, query url.Values, decodeTo any) error {
	apiURL, err := url.Parse(t.BaseURL)
	if err != nil {
		return err
	}
	query.Set("key", t.Key)
	query.Set("token", t.Token)
	apiURL = apiURL.JoinPath(path)
	apiURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error in getJSON http.NewRequestWithContext %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error in getJSON http.DefaultClient.Do %w", err)
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

func (t *trelloClient) ListBoards(ctx context.Context) ([]Board, error) {
	query := url.Values{}
	query.Set("filter", "open")
	query.Set("fields", "id,name,lists")
	query.Set("lists", "open")
	query.Set("list_fields", "id,name")

	path := "/1/members/me/boards"
	boards := make([]Board, 0)
	err := t.getJSON(ctx, path, query, &boards)
	return boards, err
}

func (t *trelloClient) ListCards(ctx context.Context, list Row) ([]Row, error) {
	query := url.Values{}
	path := "/1/lists/" + list.ID + "/cards"
	cards := make([]Row, 0)
	err := t.getJSON(ctx, path, query, &cards)
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
	return t, err
}

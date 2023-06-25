package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

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

type trelloAPI interface {
	ListBoards() ([]Board, error)
	ListCards(list Row) ([]Row, error)
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
		fmt.Fprintf(os.Stderr, "getJSON err1: %s", err)
		return err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error from Trello API: %s", string(bytes))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "getJSON err2: %s", err)
		return err
	}
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

func newTrello(baseURL string) (trelloAPI, error) {
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

func formatError(err error) string {
	message := ""
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			message += fmt.Sprintf("Invalid environment variable %s\n", strings.ToUpper(e.StructField()))
		}
	} else {
		return err.Error()
	}
	return message + "Please set your Trello API KEY and TOKEN values as environment variables.\n"
}

func run(trello trelloAPI, out io.Writer) error {
	boards, err := trello.ListBoards()
	if err != nil {
		return err
	}
	doing := regexp.MustCompile("(To Do|Doing)")
	for _, board := range boards {
		fmt.Fprintf(out, "📋%s\n", board)
		for _, list := range board.Lists {
			if !doing.MatchString(list.Name) {
				continue
			}
			fmt.Fprintf(out, "  📃%s\n", list)
			cards, err := trello.ListCards(list)
			if err != nil {
				return err
			}
			for _, card := range cards {
				fmt.Fprintf(out, "    🪧%s\n", card)
			}
		}
	}
	return nil
}

func main() {
	baseURL := "https://api.trello.com"
	trello, err := newTrello(baseURL)
	if err != nil {
		fmt.Fprint(os.Stderr, formatError(err))
		os.Exit(1)
	}
	err = run(trello, os.Stdout)
	if err != nil {
		fmt.Fprint(os.Stderr, formatError(err))
		os.Exit(1)
	}
}

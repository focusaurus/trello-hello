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

type trello struct {
	BaseURL string `validate:"min=7"`
	Key     string `validate:"min=20"`
	Token   string `validate:"min=50"`
}

func (t *trello) getJSON(path string, query url.Values, decodeTo any) error {
	query.Set("key", t.Key)
	query.Set("token", t.Token)
	res, err := http.Get(t.BaseURL + path + "?" + query.Encode())
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

func (t *trello) ListCards(list Row) ([]Row, error) {
	query := url.Values{}
	path := "/1/lists/" + list.ID + "/cards"
	cards := make([]Row, 0)
	err := t.getJSON(path, query, &cards)
	return cards, err
}

func newTrello(baseURL string) (*trello, error) {
	t := &trello{
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

func main() {
	baseURL := "https://api.trello.com"
	trello, err := newTrello(baseURL)
	if err != nil {
		fmt.Fprint(os.Stderr, formatError(err))
		os.Exit(1)
	}
	boards, err := trello.ListBoards()
	if err != nil {
		log.Fatalf("oops1: %s", err)
	}
	doing := regexp.MustCompile("(To Do|Doing)")
	for _, board := range boards {
		fmt.Printf("📋%s\n", board)
		for _, list := range board.Lists {
			if !doing.MatchString(list.Name) {
				// if strings.Contains(list.Name, "Done") {
				continue
			}
			fmt.Printf("  📃%s\n", list)
			cards, err := trello.ListCards(list)
			if err != nil {
				log.Fatalf("oops2: %s", err)
			}
			for _, card := range cards {
				fmt.Printf("    🪧%s\n", card)
			}
		}
	}
}

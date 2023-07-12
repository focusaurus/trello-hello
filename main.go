package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type trelloAPI interface {
	ListBoards(ctx context.Context) ([]Board, error)
	ListCards(ctx context.Context, list Row) ([]Row, error)
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
	ctx := context.Background()
	boards, err := trello.ListBoards(ctx)
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
			cards, err := trello.ListCards(ctx, list)
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
		os.Exit(2)
	}
}

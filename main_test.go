package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatError(t *testing.T) {
	t.Run("custom env var validation error messages", func(t *testing.T) {
		t.Setenv("KEY", "")
		t.Setenv("TOKEN", "")
		_, err := newTrello("")
		assert.Error(t, err)
		message := formatError(err)
		assert.Contains(t, message, "Invalid environment variable KEY")
		assert.Contains(t, message, "Invalid environment variable TOKEN")
		assert.Contains(t, message, "Please set your Trello API")
	})

	t.Run("error but not validation", func(t *testing.T) {
		err := errors.New("foo")
		assert.Equal(t, err.Error(), formatError(err))
	})
}

type mockTrello struct {
	listBoardsError error
	boards          []Board
	listCardsError  error
	cards           []Row
}

func (t *mockTrello) ListBoards(ctx context.Context) ([]Board, error) {
	return t.boards, t.listBoardsError
}

func (t *mockTrello) ListCards(ctx context.Context, list Row) ([]Row, error) {
	return t.cards, t.listCardsError
}

var fakeBoards = []Board{
	{Row: Row{Name: "Boardy", ID: "abcd1"}, Lists: []Row{{ID: "To Do: listb1l1", Name: "Doing: List B1L1"}}},
	{Row: Row{Name: "B2", ID: "abcd2"}, Lists: []Row{{ID: "skiplist1", Name: "Skip This List"}}},
}
var fakeCards = []Row{
	{Name: "Card 1", ID: "card1"},
	{Name: "Card 2", ID: "card2"},
}

func TestRunBaseCase(t *testing.T) {
	out := bytes.NewBufferString("")
	trello := &mockTrello{boards: fakeBoards, cards: fakeCards}
	err := run(trello, out)
	assert.NoError(t, err)
	outString := out.String()
	assert.Contains(t, outString, "📋Boardy")
	assert.Contains(t, outString, "📋B2")
	assert.Contains(t, outString, "📃Doing: List B1L1")
	assert.Contains(t, outString, "🪧Card 1")
	assert.Contains(t, outString, "🪧Card 2")
	assert.NotContains(t, outString, "Skip This List")
}

func TestRunErrors(t *testing.T) {
	t.Run("ListBoards error", func(t *testing.T) {
		trello := &mockTrello{listBoardsError: assert.AnError}
		err := run(trello, bytes.NewBufferString(""))
		assert.Error(t, err)
		assert.Equal(t, err, assert.AnError)
	})

	t.Run("ListCards error", func(t *testing.T) {
		trello := &mockTrello{boards: fakeBoards, listCardsError: assert.AnError}
		err := run(trello, bytes.NewBufferString(""))
		assert.Error(t, err)
		assert.Equal(t, err, assert.AnError)
	})
}

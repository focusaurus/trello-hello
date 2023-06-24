package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrello(t *testing.T) {
	base := "https://unittestbaseurl"
	key := "unittestkey"
	token := "unittesttoken"
	t.Setenv("KEY", key)
	t.Setenv("TOKEN", token)
	trello := newTrello(base)
	assert.Equal(t, trello.baseURL, base)
	assert.Equal(t, trello.key, key)
	assert.Equal(t, trello.token, token)
}

func TestListBoards(t *testing.T) {
	tests := map[string]struct {
		res      string
		expected []Board
	}{
		"empty response": {
			res: "[]",
			// expected: []Board{{Name: "foo", ID: "bar"}, {Name: "foo2", ID: "blah"}},
			expected: []Board{},
		},
		"base case": {
			res: `[
			{
				"name": "Boardy",
				"id": "abcd1",
				"lists": [{"id":"listb1l1", "name":"List B1L1"}]
			}, {
			"name":"B2",
			"id":"abcd2"
			}]`,
			expected: []Board{
				{Row: Row{Name: "Boardy", ID: "abcd1"}, Lists: []Row{{ID: "listb1l1", Name: "List B1L1"}}},
				{Row: Row{Name: "B2", ID: "abcd2"}},
			},
		},
	}
	for _, test := range tests {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assert.True(t, query.Has("key"))
			assert.True(t, query.Has("token"))
			assert.Equal(t, query.Get("filter"), "open")
			assert.Equal(t, query.Get("fields"), "id,name,lists")
			assert.Equal(t, query.Get("lists"), "open")
			assert.Equal(t, query.Get("list_fields"), "id,name")
			strings.HasPrefix(r.URL.Path, "/1/members/me/boards")
			responseBody := test.res
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		}))
		defer testServer.Close()
		trello := newTrello(testServer.URL)
		boards, err := trello.ListBoards()
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, boards)
	}
}

func TestListCards(t *testing.T) {
	tests := map[string]struct {
		res      string
		expected []Row
	}{
		"empty response": {
			res:      "[]",
			expected: []Row{},
		},
		"base case": {
			res: `[{"name": "Card 1", "id": "card1"}, {"name":"Card 2", "id":"card2"}]`,
			expected: []Row{
				{Name: "Card 1", ID: "card1"},
				{Name: "Card 2", ID: "card2"},
			},
		},
	}
	for _, test := range tests {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assert.True(t, query.Has("key"))
			assert.True(t, query.Has("token"))
			strings.HasPrefix(r.URL.Path, "/1/lists/list1/cards")
			responseBody := test.res
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		}))
		defer testServer.Close()
		trello := newTrello(testServer.URL)
		boards, err := trello.ListCards(Row{ID: "list1"})
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, boards)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListBoards(t *testing.T) {
	tests := map[string]struct {
		res      string
		expected []Board
	}{
		"empty response": {
			res: "[]",
			// expected: []Board{{Name: "foo", Id: "bar"}, {Name: "foo2", Id: "blah"}},
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
				{Name: "Boardy", Id: "abcd1", Lists: []List{{Id: "listb1l1", Name: "List B1L1"}}},
				{Name: "B2", Id: "abcd2"},
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

func TestListLists(t *testing.T) {
	tests := map[string]struct {
		res      string
		expected []List
	}{
		"empty response": {
			res:      "[]",
			expected: []List{},
		},
		"base case": {
			res: `[{"name": "Listy1", "id": "abcd1"}, {"name":"Listy2", "id":"abcd2"}]`,
			expected: []List{
				{Name: "Listy1", Id: "abcd1"},
				{Name: "Listy2", Id: "abcd2"},
			},
		},
	}
	for _, test := range tests {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assert.True(t, query.Has("key"))
			assert.True(t, query.Has("token"))
			strings.HasPrefix(r.URL.Path, "/1/boards/idboard1/lists")
			responseBody := test.res
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		}))
		defer testServer.Close()
		trello := newTrello(testServer.URL)
		boards, err := trello.ListLists(Board{Id: "idboard1"})
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, boards)
	}
}

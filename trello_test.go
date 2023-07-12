package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fakeKey   = "fakekey12345678901234567890"
	fakeToken = "faketoken12345678901234567890123456789012345678901234567890"
)

func TestListBoards(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)
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
			assert.True(t, strings.HasPrefix(r.URL.Path, "/1/members/me/boards"))
			responseBody := test.res
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseBody))
			assert.NoError(t, err)
		}))
		defer testServer.Close()
		trello, err := newTrello(testServer.URL)
		assert.NoError(t, err)
		boards, err := trello.ListBoards()
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, boards)
	}
}

func TestListCards(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)
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
			assert.True(t, strings.HasPrefix(r.URL.Path, "/1/lists/list1/cards"))
			responseBody := test.res
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseBody))
			assert.NoError(t, err)
		}))
		defer testServer.Close()
		trello, err := newTrello(testServer.URL)
		assert.NoError(t, err)

		boards, err := trello.ListCards(Row{ID: "list1"})
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, boards)
	}
}

func TestAPIUnauthorized(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)
	responseBody := "Response body from trello"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(responseBody))
		assert.NoError(t, err)
	}))
	defer testServer.Close()
	trello, err := newTrello(testServer.URL)
	assert.NoError(t, err)

	_, err = trello.ListCards(Row{ID: "list1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), responseBody)
}

func TestAPIInvalidJSON(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("this is not valid json"))
		assert.NoError(t, err)
	}))
	defer testServer.Close()
	trello, err := newTrello(testServer.URL)
	assert.NoError(t, err)

	_, err = trello.ListCards(Row{ID: "list1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestAPIHTTPError(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)
	trello, err := newTrello("https://nosuchhost424242")
	assert.NoError(t, err)
	_, err = trello.ListCards(Row{ID: "list1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such host")
}

func TestAPIResponseIOError(t *testing.T) {
	t.Setenv("KEY", fakeKey)
	t.Setenv("TOKEN", fakeToken)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		// Send inaccurate content length to force IO error when reading body
		w.Header().Add("Content-Length", "1024")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("[]"))
		assert.NoError(t, err)
	}))
	defer testServer.Close()
	trello, err := newTrello(testServer.URL)
	assert.NoError(t, err)

	_, err = trello.ListCards(Row{ID: "list1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected EOF")
}

func TestNewTrello(t *testing.T) {
	_, err := newTrello("imap://nothttp")
	assert.Error(t, err)
}

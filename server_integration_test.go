package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taton825/assessment/config"
	"github.com/taton825/assessment/expense"
)

func TestCreateExpense(t *testing.T) {

	config.LoadEnvironmentLocal()

	body := bytes.NewBufferString(`{
		"title": "test server integration title",
		"amount": 10,
		"note": "test server integration note", 
		"tags": ["integration", "test"]
	}`)
	var e expense.Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0.0, e.ID)
	assert.Equal(t, "test server integration title", e.Title)
	assert.Equal(t, 10.0, e.Amount)
	assert.Equal(t, "test server integration note", e.Note)
	assert.Greater(t, len(e.Tags), 0)
	assert.Equal(t, "integration", e.Tags[0])
	assert.Equal(t, "test", e.Tags[1])
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

func TestGetExpense(t *testing.T) {

	config.LoadEnvironmentLocal()

	e := createExpense(t)
	fmt.Println(e.ID)
	var latest expense.Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latest.ID)
	assert.Equal(t, e.Title, latest.Title)
	assert.Equal(t, e.Amount, latest.Amount)
	assert.Equal(t, e.Note, latest.Note)
	assert.Equal(t, e.Tags, latest.Tags)
}

func createExpense(t *testing.T) expense.Expense {
	var e expense.Expense
	body := bytes.NewBufferString(`{
		"title": "test server integration title by test",
		"amount": 20,
		"note": "test server integration note by test", 
		"tags": ["integration by test", "tag by test"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("Can't create expense: ", err)
	}
	return e
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
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

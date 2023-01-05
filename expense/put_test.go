//go:build unit
// +build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPutExpenseHandler(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(`{
		"title": "test handler title",
		"amount": 10,
		"note": "test handler note", 
		"tags": ["handler", "test"]
	}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tags := []string{"handler", "test"}
	updateMockSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id"

	updateMockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")

	db, mock, err := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta(updateMockSql)).
		WithArgs("1", "test handler title", 10.0, "test handler note", pq.Array(tags)).
		WillReturnRows(updateMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Act
	err = h.PutExpenseHandler(c)

	// Assertion
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

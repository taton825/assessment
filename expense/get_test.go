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

func TestGetExpenseHandler(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tags := []string{"handler", "test"}
	expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test handler title", 10.0, "test handler note", pq.Array(tags))
	queryMockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(queryMockSql)).
		ExpectQuery().
		WithArgs("1").WillReturnRows(expenseMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expect := `{"id":1,"title":"test handler title","amount":10,"note":"test handler note","tags":["handler","test"]}`

	// Act
	err = h.GetExpenseHandler(c)

	// Assertion
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expect, strings.TrimSpace(rec.Body.String()))
	}
}

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

func TestCreateExpenseHandler(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(`{
		"title": "test handler title",
		"amount": 10,
		"note": "test handler note", 
		"tags": ["handler", "test"]
	}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tags := []string{"handler", "test"}
	mockedSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id"

	mockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")

	db, mock, err := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta(mockedSql)).
		WithArgs("test handler title", 10.0, "test handler note", pq.Array(tags)).
		WillReturnRows(mockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)

	// Act
	err = h.CreateExpenseHandler(c)

	// Assertion
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

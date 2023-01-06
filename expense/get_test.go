package expense

import (
	"fmt"
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

	t.Run("Test Get One Row Success", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

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
		assertResult(t, err, rec, expect)
	})

	t.Run("Test Get One Row Prepare Error", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

		tags := []string{"handler", "test"}
		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test handler title", 10.0, "test handler note", pq.Array(tags))
		queryMockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1$" // Error Here

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

		// Act
		err = h.GetExpenseHandler(c)

		// Assertion
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Test Get One Row Scan Error", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

		tags := []string{"handler", "test"}
		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test handler title", 10.0, "test handler note", pq.Array(tags))
		queryMockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"

		db, mock, err := sqlmock.New()
		mock.ExpectPrepare(regexp.QuoteMeta(queryMockSql)).
			ExpectQuery().
			WithArgs(1).WillReturnRows(expenseMockRows) // Error Here
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Act
		err = h.GetExpenseHandler(c)

		// Assertion
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetExpensesHandler(t *testing.T) {

	t.Run("Test Get All Row Success", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

		tags := []string{"handler", "test"}

		expensesMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test handler title1", 10.0, "test handler note1", pq.Array(tags)).
			AddRow("2", "test handler title2", 10.0, "test handler note2", pq.Array(tags))
		queryMockSql := "SELECT id, title, amount, note, tags FROM expenses"

		db, mock, err := sqlmock.New()
		mock.ExpectPrepare(regexp.QuoteMeta(queryMockSql)).
			ExpectQuery().
			WillReturnRows(expensesMockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}
		c := e.NewContext(req, rec)

		// Act
		err = h.GetExpensesHandler(c)

		expect := `[{"id":1,"title":"test handler title1","amount":10,"note":"test handler note1","tags":["handler","test"]},{"id":2,"title":"test handler title2","amount":10,"note":"test handler note2","tags":["handler","test"]}]`

		// Assertion
		assertResult(t, err, rec, expect)
	})

	t.Run("Test Get All Row Prepare Error", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

		tags := []string{"handler", "test"}

		expensesMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test handler title1", 10.0, "test handler note1", pq.Array(tags)).
			AddRow("2", "test handler title2", 10.0, "test handler note2", pq.Array(tags))
		queryMockSql := "SELECT id, title, amount, note, tags FROM expenses $"

		db, mock, err := sqlmock.New()
		mock.ExpectPrepare(regexp.QuoteMeta(queryMockSql)).
			ExpectQuery().
			WillReturnRows(expensesMockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}
		c := e.NewContext(req, rec)

		// Act
		err = h.GetExpensesHandler(c)

		// Assertion
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Test Get All Row Query Error", func(t *testing.T) {
		// Arrange
		e, req, rec := initialRequest()

		tags := []string{"handler", "test"}

		expensesMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "test handler title1", 10.0, "test handler note1", pq.Array(tags)).
			AddRow("2", "test handler title2", "10.0", "test handler note2", pq.Array(tags))

		queryMockSql := "SELECT id, title, amount, note, tags FROM expenses"

		db, mock, err := sqlmock.New()
		mock.ExpectPrepare(regexp.QuoteMeta(queryMockSql)).
			ExpectQuery().
			WillReturnRows(expensesMockRows).
			WillReturnError(fmt.Errorf("Error Query"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}
		c := e.NewContext(req, rec)

		// Act
		err = h.GetExpensesHandler(c)

		// Assertion
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func initialRequest() (*echo.Echo, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e, req, rec
}

func assertResult(t *testing.T, err error, rec *httptest.ResponseRecorder, expect string) {
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expect, strings.TrimSpace(rec.Body.String()))
	}
}

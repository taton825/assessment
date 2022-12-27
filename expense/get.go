package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/taton825/assessment/database"
)

func GetExpenseHandler(c echo.Context) error {
	rowid := c.Param("id")

	stmt, err := database.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't create prepare statement for get one expense: " + err.Error()})
	}

	row := stmt.QueryRow(rowid)

	var e Expense
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan row into variable: " + err.Error()})
	}
	return c.JSON(http.StatusOK, e)
}

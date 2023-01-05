package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/taton825/assessment/database"
)

func (h *handler) GetExpenseHandler(c echo.Context) error {
	rowid := c.Param("id")

	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	// fmt.Println("err prepare: " + err.Error())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't create prepare statement for get one expense: " + err.Error()})
	}

	row := stmt.QueryRow(rowid)

	var e Expense
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	// fmt.Println("err scan: " + err.Error())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan row into variable: " + err.Error()})
	}
	return c.JSON(http.StatusOK, e)
}

func GetExpensesHandler(c echo.Context) error {
	stmt, err := database.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't create prepare statement for get one expense: " + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses:" + err.Error()})
	}

	var expenses []Expense

	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan row into variable:" + err.Error()})
		}
		expenses = append(expenses, e)
	}
	log.Println("query all expenses success")

	return c.JSON(http.StatusOK, expenses)
}

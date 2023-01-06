package expense

import (
	"database/sql"
	"testing"
)

func TestNewApplication(t *testing.T) {
	db, _ := sql.Open("mock", "")
	app := NewApplication(db)

	if app.DB != db {
		t.Errorf("Expected app.DB to be %v, but got %v", db, app.DB)
	}
}

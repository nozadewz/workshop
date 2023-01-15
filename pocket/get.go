package pocket

import ("github.com/labstack/echo/v4"
"database/sql"
"net/http"
)
func (h handler) GetPocketBalanceById(c echo.Context) error {
	id := c.Param("id")
	ex := Pocket{}

	row := h.db.QueryRow("SELECT * FROM pocket WHERE id = $1", id)
	err := row.Scan(&ex.ID, &ex.Balance, &ex.Currency)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "pocket not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't find pocket please contact admin:" + err.Error()})
	}
}
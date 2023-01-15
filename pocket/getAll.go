package pocket

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) GetAllPockets(c echo.Context) error {

	id := c.Param("id")
	pockets, err := getAllPockets(h.db, id)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pockets)
}

func getAllPockets(db *sql.DB, id string) ([]Pocket, error) {

	queryStatement := `
	SELECT * FROM pockets
	`
	st, err := db.Prepare(queryStatement)

	if err != nil {
		return nil, err
	}

	rows, err2 := st.Query()
	if err2 != nil {
		return nil, err
	}

	pockets := []Pocket{}

	for rows.Next() {

		pocket := Pocket{}

		err = rows.Scan(&pocket.ID, &pocket.Account_ID, &pocket.Name, &pocket.Category, &pocket.Currency, &pocket.Balance)

		if err != nil {
			return pockets, err
		}

		pockets = append(pockets, pocket)
	}

	return pockets, nil

}

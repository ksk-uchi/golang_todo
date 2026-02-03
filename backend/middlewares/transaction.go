package middlewares

import (
	"todo-app/ent"

	"github.com/labstack/echo/v5"
)

func NewTransactionMiddleware(client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tx, err := client.Tx(c.Request().Context())
			if err != nil {
				return err
			}

			// Context にセット
			c.Set("tx", tx)

			if err := next(c); err != nil {
				tx.Rollback()
				return err
			}
			return tx.Commit()
		}
	}
}

func GetTx(c *echo.Context) *ent.Tx {
	return c.Get("tx").(*ent.Tx)
}

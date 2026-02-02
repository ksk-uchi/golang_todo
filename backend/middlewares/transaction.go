package middlewares

import (
	"todo-app/ent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v5"
)

func TransactionMiddleware(client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tx, err := client.Tx(c.Request().Context())
			if err != nil {
				return err
			}

			// トランザクションを Context にセット
			c.Set("tx", tx)

			// 次の処理（ハンドラー）を実行
			if err := next(c); err != nil {
				tx.Rollback() // エラーがあればロールバック
				return err
			}

			return tx.Commit() // 正常終了ならコミット
		}
	}
}

func GetTx(c *echo.Context) *ent.Tx {
	return c.Get("tx").(*ent.Tx)
}

package repositories

import (
	"context"
	"todo-app/ent"
)

func getEntClient(ctx context.Context) *ent.Client {
	tx := ent.TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return ent.FromContext(ctx)
}

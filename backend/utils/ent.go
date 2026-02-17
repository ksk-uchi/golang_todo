package utils

import (
	"context"
	"todo-app/ent"
)

// WithTx creates a new transaction and returns a new context with the transaction attached.
// It uses ent.NewTxContext to attach the transaction to the context.
func WithTx(ctx context.Context, client *ent.Client) (context.Context, *ent.Tx, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	// ent.NewTxContext attaches the transaction to the context,
	// allowing ent.TxFromContext to retrieve it.
	return ent.NewTxContext(ctx, tx), tx, nil
}

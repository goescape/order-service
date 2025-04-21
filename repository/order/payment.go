package repository

import (
	"context"
	"log"
	"order-svc/model"
)

func (o *orderStore) PayOrder(ctx context.Context, req *model.PayOrderModel) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		log.Default().Println("Failed to begin transaction:", err)
		return err
	}
	defer tx.Rollback()

	var orderId string

	query := `
		UPDATE orders set status = 'paid'
		where id = $1 AND status = 'pending'
		returning id
	`

	err = tx.QueryRowContext(ctx, query, req.ID).Scan(&orderId)
	if err != nil {
		log.Default().Println("Failed to paid order:", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}

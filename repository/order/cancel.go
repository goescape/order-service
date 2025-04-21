package repository

import (
	"context"
	"log"
	"order-svc/model"
)

func (o *orderStore) CancelOrder(ctx context.Context, req *model.CancelOrderModel) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		log.Default().Println("Failed to begin transaction:", err)
		return err
	}
	defer tx.Rollback()

	var orderId string

	query := `
		UPDATE orders set status = 'cancel'
		where id = $1 AND status = 'pending'
		returning id
	`

	err = tx.QueryRowContext(ctx, query, req.ID).Scan(&orderId)
	if err != nil {
		log.Default().Println("Failed to cancel order:", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}

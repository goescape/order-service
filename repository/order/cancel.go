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

	type OrderDetail struct {
		ProductId string
		Qty       int
	}

	var orderDetails []OrderDetail

	queryOrderDetails := `
		SELECT product_id, qty FROM order_details
		WHERE order_id = $1
	`

	rows, err := tx.QueryContext(ctx, queryOrderDetails, orderId)
	if err != nil {
		log.Default().Println("Failed to query order details:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var detail OrderDetail
		if err := rows.Scan(&detail.ProductId, &detail.Qty); err != nil {
			log.Default().Println("Failed to scan order detail:", err)
			return err
		}

		orderDetails = append(orderDetails, detail)
	}
	if err := rows.Err(); err != nil {
		log.Default().Println("Error iterating over order details:", err)
		return err
	}

	// Update the stock for each product in the order details
	for _, detail := range orderDetails {
		queryUpdateStock := `
			UPDATE products
			SET qty = qty + $1
			WHERE id = $2
		`

		_, err := tx.ExecContext(ctx, queryUpdateStock, detail.Qty, detail.ProductId)
		if err != nil {
			log.Default().Println("Failed to update stock:", err)
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}

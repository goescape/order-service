package repository

import (
	"context"
	"log"
	"order-svc/helpers/dbutil"
	"order-svc/model"

	"github.com/shopspring/decimal"
)

func (o *orderStore) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		log.Default().Println("Failed to begin transaction:", err)
		return nil, err
	}
	defer tx.Rollback()

	var orderId string

	query := `
		INSERT INTO orders (user_id, total_price)
		VALUES ($1, $2)
		RETURNING id
	`

	// counting total price
	var totalPrice decimal.Decimal

	for _, item := range req.Items {
		price := decimal.NewFromFloat(item.Price)
		qty := decimal.NewFromInt(item.Qty)

		totalPrice = totalPrice.Add(price.Mul(qty))
	}

	err = tx.QueryRowContext(ctx, query, req.UserId, totalPrice).Scan(&orderId)
	if err != nil {
		log.Default().Println("Failed to insert order:", err)
		return nil, err
	}

	// Insert order details
	queryDetail := `
		INSERT INTO order_details (order_id, product_id, price, qty)
		VALUES
	`
	args := make([]interface{}, 0, len(req.Items)*4)

	for i, item := range req.Items {
		if i > 0 {
			queryDetail += ", "
		}
		queryDetail += `(?, ?, ?, ?)`
		args = append(args, orderId, item.ProductId, item.Price, item.Qty)
	}

	_, err = tx.ExecContext(ctx, dbutil.ReplacePlaceholders(queryDetail), args...)
	if err != nil {
		log.Default().Println("Failed to insert order details:", err)
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &model.CreateOrderResp{
		OrderId: orderId,
	}, nil
}

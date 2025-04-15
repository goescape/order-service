package repository

import (
	"context"
	"log"
	"order-svc/model"
)

const (
	GetOrderList = `
	SELECT
	COUNT(*) OVER() AS total_data,
    o.id AS order_id,
    o.user_id AS user_id,
    o.total_price AS total_price,
    o.created_at AS created_at,
    o.updated_at AS updated_at,
    od.id AS detail_id,
    od.product_id AS product_id,
    od.qty AS qty,
    od.price AS price,
	FROM orders o
	LEFT JOIN order_details od ON o.id = od.order_id
	where o.user_id = $1
	ORDER BY o.created_at DESC
	LIMIT $2 OFFSET $3;
	`
)

func (s *orderStore) GetOrderList(ctx context.Context, req *model.GetOrderListRequest) (*model.ListOrderResponse, error) {

	var (
		totalData uint32
		resp      = new(model.ListOrderResponse)
	)

	rows, err := s.db.QueryContext(ctx, GetOrderList, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		log.Default().Println("failed to query GetOrderList:", err)
		return nil, err
	}
	defer rows.Close()

	var orderData []*model.OrderModel
	for rows.Next() {
		var data model.OrderModel

		if err := rows.Scan(&totalData, &data.ID, &data.UserID, &data.TotalPrice, &data.CreatedAt, &data.UpdatedAt,
			&data.DetailID, &data.ProductID, &data.Qty, &data.Price); err != nil {
			log.Default().Println("failed to scan orderlist:", err)
			return nil, err
		}
		orderData = append(orderData, &data)
	}
	if err := rows.Err(); err != nil {
		log.Default().Println("failed to iterate orderList:", err)
		return nil, err
	}

	OrderListData := model.MapOrderModelsToResponse(orderData)
	resp.Items = OrderListData

	return resp, nil
}

package model

import "time"

type OrderDetailResponse struct {
	ID        string  `json:"detail_id" db:"order_detail_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	Qty       int     `json:"qty" db:"qty"`
	Price     float64 `json:"price" db:"price"`
}

type OrderResponse struct {
	ID         string                `json:"order_id" db:"order_id"`
	UserID     string                `json:"user_id" db:"user_id"`
	TotalPrice float64               `json:"total_price" db:"total_price"`
	CreatedAt  time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at" db:"updated_at"`
	Details    []OrderDetailResponse `json:"details"` // biasanya bagian array ini diisi manual setelah scan
}

type GetOrderListRequest struct {
	UserId *string `json:"user_id"`
	Page   uint32  `json:"page,omitempty"`
	Limit  uint32  `json:"limit,omitempty"`
}

type OrderModel struct {
	ID         string    `json:"order_id" db:"order_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DetailID   string    `json:"detail_id" db:"order_detail_id"`
	ProductID  string    `json:"product_id" db:"product_id"`
	Qty        int       `json:"qty" db:"qty"`
	Price      float64   `json:"price" db:"price"`
}

type ListOrderResponse struct {
	Items []*OrderResponse `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Meta  *Meta            `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
}

type Meta struct {
	TotalData   uint32 `protobuf:"varint,1,opt,name=total_data,json=totalData,proto3" json:"total_data,omitempty"`
	TotalPage   uint32 `protobuf:"varint,2,opt,name=total_page,json=totalPage,proto3" json:"total_page,omitempty"`
	CurrentPage uint32 `protobuf:"varint,3,opt,name=current_page,json=currentPage,proto3" json:"current_page,omitempty"`
	Limit       uint32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func MapOrderModelsToResponse(models []*OrderModel) []*OrderResponse {
	orderMap := make(map[string]*OrderResponse)

	for _, m := range models {
		// Kalau order belum ada, buat baru
		if _, exists := orderMap[m.ID]; !exists {
			orderMap[m.ID] = &OrderResponse{
				ID:         m.ID,
				UserID:     m.UserID,
				TotalPrice: m.TotalPrice,
				CreatedAt:  m.CreatedAt,
				UpdatedAt:  m.UpdatedAt,
				Details:    []OrderDetailResponse{},
			}
		}

		// Kalau DetailID kosong, artinya tidak ada detail (LEFT JOIN mungkin null)
		if m.DetailID != "" {
			detail := OrderDetailResponse{
				ID:        m.DetailID,
				ProductID: m.ProductID,
				Qty:       m.Qty,
				Price:     m.Price,
			}
			orderMap[m.ID].Details = append(orderMap[m.ID].Details, detail)
		}
	}

	// Ubah map ke slice
	responses := make([]*OrderResponse, 0, len(orderMap))
	for _, resp := range orderMap {
		responses = append(responses, resp)
	}

	return responses
}

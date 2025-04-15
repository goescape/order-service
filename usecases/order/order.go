package usecases

import (
	"context"
	"order-svc/model"
	repository "order-svc/repository/order"
)

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
}

type orderUsecase struct {
	repo repository.OrderRepository
}

var _ OrderUsecases = &orderUsecase{}

func NewOrderUsecase(repo repository.OrderRepository) *orderUsecase {
	return &orderUsecase{
		repo: repo,
	}
}

func (o *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	resp, err := o.repo.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

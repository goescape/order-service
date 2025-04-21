package usecases

import (
	"context"
	"log"
	"order-svc/model"
	repository "order-svc/repository/order"
	rdscheduler "order-svc/repository/redis/scheduler"
)

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	GetOrderList(ctx context.Context, req *model.GetOrderListRequest) (*model.ListOrderResponse, error)
	PayOrder(ctx context.Context, req *model.PayOrderModel) error
}

type orderUsecase struct {
	repo        repository.OrderRepository
	rdscheduler rdscheduler.SchedulerInterface
}

var _ OrderUsecases = &orderUsecase{}

func NewOrderUsecase(repo repository.OrderRepository, rd rdscheduler.SchedulerInterface) *orderUsecase {
	return &orderUsecase{
		repo:        repo,
		rdscheduler: rd,
	}
}

func (o *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	resp, err := o.repo.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	err = o.rdscheduler.ScheduleTaskCancellation(resp.OrderId)
	if err != nil {
		log.Println(err)
	}

	return resp, nil
}

func (o *orderUsecase) GetOrderList(ctx context.Context, req *model.GetOrderListRequest) (*model.ListOrderResponse, error) {
	resp, err := o.repo.GetOrderList(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *orderUsecase) PayOrder(ctx context.Context, req *model.PayOrderModel) error {
	err := o.repo.PayOrder(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

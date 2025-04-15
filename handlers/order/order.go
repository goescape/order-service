package handlers

import (
	"fmt"
	"log"
	"order-svc/helpers/fault"
	"order-svc/helpers/response"
	"order-svc/model"
	usecases "order-svc/usecases/order"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	createOrderLock sync.Mutex
)

type OrderHandler struct {
	service usecases.OrderUsecases
}

func NewOrderHandler(s usecases.OrderUsecases) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		fault.Response(c, fault.Custom(
			400,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	if len(req.Items) == 0 {
		log.Default().Println("Items cannot be empty")
		fault.Response(c, fault.Custom(
			400,
			fault.ErrBadRequest,
			"items cannot be empty",
		))
		return
	}

	resp, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		fault.Response(c, err)
		return
	}

	c.JSON(200, resp)
}

func (h *OrderHandler) GetOrderList(c *gin.Context) {
	var req model.GetOrderListRequest

	userId := c.DefaultQuery("user_id", "")
	if userId != "" {
		req.UserId = &userId
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	resp, err := h.service.GetOrderList(c.Request.Context(), &req)
	if err != nil {
		fault.Response(c, err)
		return
	}

	response.JSON(c, 200, "Success", resp)
}

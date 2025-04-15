package handlers

import (
	"fmt"
	"order-svc/helpers/fault"
	"order-svc/helpers/response"
	"order-svc/model"
	usecases "order-svc/usecases/order"

	"github.com/gin-gonic/gin"
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

	resp, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		fault.Response(c, err)
		return
	}

	response.JSON(c, 200, "Success", resp)
}

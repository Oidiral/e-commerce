package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/service"
)

type CartHandler struct {
	svc service.CartService
}

func NewCartHandler(svc service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
	Qty       int       `json:"qty"`
}

type ChangeQtyRequest struct {
	Qty int `json:"qty"`
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	cart, err := h.svc.GetCart(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	if err := h.svc.AddItem(c.Request.Context(), userID, req.ProductID, req.Qty); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *CartHandler) ChangeQty(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	var req ChangeQtyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	if err := h.svc.ChangeQty(c.Request.Context(), userID, productID, req.Qty); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	if err := h.svc.RemoveItem(c.Request.Context(), userID, productID); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *CartHandler) Clear(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		HandleError(c, service.ErrBadRequest)
		return
	}
	if err := h.svc.Clear(c.Request.Context(), userID); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

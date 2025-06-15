package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/service"
)

func RegisterRoutes(router *gin.Engine, svc service.CartService) {
	h := NewCartHandler(svc)
	api := router.Group("/api/v1/cart")
	{
		api.GET("/:id", h.GetCart)
		api.POST("/:id/items", h.AddItem)
		api.PUT("/:id/items/:product_id", h.ChangeQty)
		api.DELETE("/:id/items/:product_id", h.RemoveItem)
		api.DELETE("/:id/items", h.Clear)
	}
}

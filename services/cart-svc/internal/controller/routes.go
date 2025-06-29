package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/service"
)

func RegisterRoutes(router *gin.Engine, svc service.CartService) {
	h := NewCartHandler(svc)
	api := router.Group("/api/v1/cart")
	{
		api.GET("/:user_id", h.GetCart)
		api.POST("/:user_id/items", h.AddItem)
		api.PUT("/:user_id/items/:product_id", h.ChangeQty)
		api.DELETE("/:user_id/items/:product_id", h.RemoveItem)
		api.DELETE("/:user_id/items", h.Clear)
	}
}

package http

import (
	"github.com/labstack/echo/v4"
	"github.com/notblessy/handler"
)

type HTTPService struct {
	handler *handler.Handler
}

// New :nodoc:
func New(h *handler.Handler) *HTTPService {
	return &HTTPService{
		handler: h,
	}
}

func (h *HTTPService) Routes(route *echo.Echo) {
	v1 := route.Group("/v1")
	v1.POST("/splits/:owner_id", h.handler.SplitBill)
	v1.GET("/splits/:id", h.handler.FindSplit)

	bill := v1.Group("/bills")
	bill.GET("", h.handler.FindAllBill)
	bill.POST("", h.handler.CreateBill)
	bill.DELETE("/:id", h.handler.DeleteBill)

	share := v1.Group("/shares")
	share.POST("", h.handler.CreateShare)
	share.GET("", h.handler.FindAllWithShare)

	user := v1.Group("/users")
	user.POST("", h.handler.CreateUser)
	user.POST("/mates", h.handler.CreateMates)
	user.GET("/:id", h.handler.FindOneUser)
	user.DELETE("/:mate_id/mate", h.handler.DeleteMate)
}

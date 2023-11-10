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

	bill := v1.Group("/bills")
	bill.GET("", h.handler.FindAllBill)
	bill.POST("", h.handler.CreateBill)
	bill.DELETE("/:id", h.handler.DeleteBill)
}

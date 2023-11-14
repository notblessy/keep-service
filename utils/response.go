package utils

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// DefaultMessage :nodoc:
const DefaultMessage string = "success"

// HTTPResponse :nodoc:
type HTTPResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Response :nodoc:
func Response(c echo.Context, status int, response *HTTPResponse) error {
	if response.Message == "" {
		response.Message = DefaultMessage
	}

	return c.JSON(status, response)
}

func Dump(v interface{}) string {
	js, _ := json.Marshal(v)
	return string(js)
}

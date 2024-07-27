package helper

import (
	"github.com/labstack/echo/v4"
)

type (
	BaseResponse struct {
		Messages string      `json:"messages"`
		Data     interface{} `json:"data"`
		Error    error       `json:"error"`
		Meta     *Pagination `json:"meta,omitempty"`
	}
)

// NewResponses return dynamic JSON responses
func NewResponses[T any](ctx echo.Context, statusCode int, message string, data T, err error, meta *Pages) error {
	if statusCode < 400 {
		if meta != nil {
			return ctx.JSON(statusCode, &BaseResponse{
				Messages: message,
				Data:     data,
				Error:    nil,
				Meta:     &meta.Pagination,
			})
		}

		return ctx.JSON(statusCode, &BaseResponse{
			Messages: message,
			Data:     data,
			Error:    nil,
		})

	}

	return ctx.JSON(statusCode, &BaseResponse{
		Messages: message,
		Data:     data,
		Error:    err,
	})
}

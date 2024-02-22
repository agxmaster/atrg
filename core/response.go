package core

import (
	"github.com/agxmaster/atrg/adapter"
	"net/http"
)

type BusinessResponse struct {
	Code    ResponseCode `json:"code"`
	Data    interface{}  `json:"data"`
	Message string       `json:"message"`
}

func Success(c adapter.GinCtxCore, data interface{}) {
	c.JSON(http.StatusOK, BusinessResponse{Code: Mp.SuccessCode, Data: data})
}

func Error(c adapter.GinCtxCore, err error) {
	c.JSON(http.StatusOK, BusinessResponse{Code: Mp.ErrorCode, Message: err.Error()})
}

package response

import (
	"net/http"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 20000,
		Msg:  "OK",
	})
}

func SuccessWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 20000,
		Msg:  msg,
	})
}

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: 20000,
		Msg:  "OK",
		Data: data,
	})
}

func Error(c *gin.Context, err error) {
	if bizErr, ok := gerrors.FromBizStatusError(err); ok {
		c.JSON(http.StatusOK, Response{
			Code: bizErr.BizStatusCode(),
			Msg:  bizErr.BizMessage(),
			Data: nil,
		})
	}

	c.JSON(http.StatusOK, Response{
		Code: 50000,
		Msg:  err.Error(),
		Data: nil,
	})
}

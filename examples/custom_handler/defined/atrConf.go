package defined

import (
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/handler"
	"github.com/gin-gonic/gin"
)

type Custom struct {
	*handler.Atr
}

func (a Custom) List(c *gin.Context) {
	core.Success(c, "Custom")
}

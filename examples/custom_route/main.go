package main

import (
	"github.com/agxmaster/atrg"
	"github.com/agxmaster/atrg/examples/custom_route/dal"
	"github.com/agxmaster/atrg/examples/custom_route/defined"
	"github.com/gin-gonic/gin"
)

func main() {

	g := gin.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, g, atrg.WithCustomRoute(defined.ConfigCustomModel))

	g.Run(":8888")
}

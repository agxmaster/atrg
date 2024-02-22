package main

import (
	"github.com/agxmaster/atrg"
	"github.com/agxmaster/atrg/examples/custom_handler/dal"
	"github.com/agxmaster/atrg/examples/custom_handler/defined"
	"github.com/gin-gonic/gin"
)

func main() {

	g := gin.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, g, atrg.WithCustomHandler(defined.Custom{}))

	g.Run(":8888")

}

package main

import (
	"github.com/agxmaster/atrg"
	"github.com/agxmaster/atrg/examples/simple/dal"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, g)
	g.Run(":8888")

}

package main

import (
	"github.com/agxmaster/atrg"
	"github.com/agxmaster/atrg/examples/bind_model/dal"
	"github.com/agxmaster/atrg/examples/bind_model/defined"
	"github.com/gin-gonic/gin"
)

func main() {

	g := gin.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, g, atrg.WithModelConfig(defined.ConfigModelMap))

	g.Run(":8888")
}

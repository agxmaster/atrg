package main

import (
	"github.com/agxmaster/atrg"
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/examples/standard/dal"
	"github.com/agxmaster/atrg/examples/standard/defined"
	"github.com/gin-gonic/gin"
)

func main() {

	g := gin.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, g, atrg.WithModelConfig(defined.ConfigModelMap),
		atrg.WithRules(core.Rules{
			{Table: []string{"student"}, RouteTypes: []core.RouteType{core.RouteTypeList}, RuleType: core.RuleTypeAllow},
			{Table: []string{"*"}, RouteTypes: []core.RouteType{"*"}, RuleType: core.RuleTypeDeny},
		}))

	g.Run(":8888")
}

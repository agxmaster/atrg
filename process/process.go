package process

import (
	"context"
	"github.com/agxmaster/atm/clause"
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/parse"
	"github.com/gin-gonic/gin"
)

type Process interface {
	Info(id int) (interface{}, error)
	List(params map[string]interface{}, page int) (interface{}, error)
	Delete(id int) error
	Create() error
	BatchCreate() error
	Update(id int) error
}

func ProcessFactory(ctx context.Context, c *gin.Context, modelName string) Process {

	var modelConf = (&core.MpModel{}).GetMpModel(core.Method(c.Request.Method), c.Request.URL.Path, modelName)

	if modelConf == nil {
		modelConf = &core.MpModel{}
	}
	if modelConf.Model != nil {
		return &AtmWithModel{BaseStore{ModelConf: modelConf, ModelName: modelName, Ctx: ctx, C: c}}
	}
	return &AtmWithOutModel{BaseStore{ModelConf: modelConf, ModelName: modelName, Ctx: ctx, C: c}}
}

type BaseStore struct {
	ModelName string
	ModelConf *core.MpModel
	Ctx       context.Context
	C         *gin.Context
}

func (b BaseStore) ProcessList(params map[string]interface{}, page int) (clause.Clauses, error) {
	params = b.ModelConf.QueryChange(b.Ctx, params)
	return (&parse.ClauseParse{ModelType: b.ModelConf}).Parse(&params, page).GetClause()
}

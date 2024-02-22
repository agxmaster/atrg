package process

import (
	"encoding/json"
	"github.com/agxmaster/atm"
	"github.com/agxmaster/atm/clause"
	"github.com/agxmaster/atrg/core"
	"io"
)

type AtmWithOutModel struct {
	BaseStore
}

func (p *AtmWithOutModel) Info(id int) (interface{}, error) {
	db := atm.M{DB: core.Db}
	data, err := db.First(p.Ctx, p.ModelName, int64(id))
	if err != nil {
		return data, nil
	}

	res, err := p.ModelConf.FormatLine(p.Ctx, data)
	return res, err
}

func (p *AtmWithOutModel) List(params map[string]interface{}, page int) (interface{}, error) {
	db := atm.M{DB: core.Db}
	claus, err := p.BaseStore.ProcessList(params, page)

	if err != nil {
		return nil, err
	}
	data, err := db.QueryPage(p.Ctx, p.ModelName, claus)
	if err != nil {
		return nil, err
	}
	data.Data, err = p.ModelConf.FormatRowsMapList(p.Ctx, data.Data)
	return data, err
}

func (p *AtmWithOutModel) Delete(id int) error {
	return (&atm.M{DB: core.Db}).Delete(p.Ctx, p.ModelName, int64(id))
}

func (p *AtmWithOutModel) Create() error {
	db := atm.M{DB: core.Db}
	var data atm.RowsMap

	err := p.C.BindJSON(&data)

	if err != nil {
		return err
	}
	if p.ModelConf.CreateParamsHandler != nil {
		data, err = p.ModelConf.CreateParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}
	return db.Create(p.Ctx, p.ModelName, data)
}

func (p *AtmWithOutModel) BatchCreate() error {

	var data []atm.RowsMap
	body, err := io.ReadAll(p.C.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	if p.ModelConf.CreateBatchParamsHandler != nil {
		data, err = p.ModelConf.CreateBatchParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}
	return (&atm.M{DB: core.Db}).BatchCreate(p.Ctx, p.ModelName, data)
}

func (p *AtmWithOutModel) Update(id int) error {

	var data clause.ColumnMap
	body, err := io.ReadAll(p.C.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if p.ModelConf.UpdateParamsHandler != nil {
		data, err = p.ModelConf.UpdateParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}
	return (&atm.M{DB: core.Db}).Update(p.Ctx, p.ModelName, int64(id), data)
}
